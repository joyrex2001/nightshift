package agent

import (
	"github.com/golang/glog"
)

// StartTrigger will consume the triggerqueue channel and execute each
// triggers sequentially. It will block until the channel is closed.
func (a *worker) StartTrigger() {
	for trgr := range a.trigqueue {
		if err := a.triggers[trgr].Execute(); err != nil {
			glog.Errorf("Error execute trigger: %s", err)
		}
	}
}

// StopTrigger will stop the scaling loop.
func (a *worker) StopTrigger() {
	close(a.trigqueue)
}

// queueTriggers will enqueue the collected triggers as specified in the
// prodived list of trigger id's. Each trigger will be enqueued just once.
func (a *worker) queueTriggers(trgrs []string) {
	done := map[string]bool{}
	for _, trgr := range trgrs {
		if done[trgr] {
			continue
		}
		_, ok := a.triggers[trgr]
		if ok {
			a.queueTrigger(trgr)
		} else {
			glog.Errorf("Error execute trigger: invalid trigger %s", trgr)
		}
		done[trgr] = true
	}
}

// queueTrigger will add a trigger to the triggerqueue.
func (a *worker) queueTrigger(trgr string) {
	a.trigqueue <- trgr
}
