package agent

import (
	"time"

	"github.com/golang/glog"
)

type delay struct {
	event []*event
	cache map[string]delayCache
}

type delayCache struct {
	expires time.Time
	result  bool
}

const delayCheckInterval = 5 * time.Minute

// InitDelay will initialize the delay checking cache.
func (a *worker) InitDelay() {
	a.delay.event = []*event{}
	a.delay.cache = map[string]delayCache{}
}

// shouldDelayScale will lazy check given KeepAlive webhooks
// in order to see if the scaling should be postponed a bit.
func (a *worker) shouldDelayScale(keepalives []string) bool {
	for _, id := range keepalives {
		c, ok := a.delay.cache[id]
		if !ok || time.Now().After(c.expires) {
			ka, ok := a.keepalives[id]
			if !ok {
				glog.Errorf("non existing keepalive %s, ignored", id)
				continue
			}
			res := true
			if err := ka.Execute(); err != nil {
				res = false
				glog.Infof("keepalive hook %s returned: %s, will not delay scale", id, err)
			}
			c = delayCache{time.Now().Add(delayCheckInterval), res}
			a.delay.cache[id] = c
		}
		if c.result {
			return true
		}
	}
	return false
}

// delayScale will add the given event to the postponed event
// list.
func (a *worker) delayScale(e *event) {
	a.delay.event = append(a.delay.event, e)
}

// getDelayedEvents will return the current list of delayed
// events.
func (a *worker) getDelayedEvents() []*event {
	// clear list and return all currently delayed events
	ev := a.delay.event
	a.delay.event = []*event{}
	return ev
}
