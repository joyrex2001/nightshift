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
