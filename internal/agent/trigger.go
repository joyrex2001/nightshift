package agent

import (
	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type triggr struct {
	id      string
	objects []*scanner.Object
}

// StartTrigger will consume the triggerqueue channel and execute each
// triggers sequentially. It will block until the channel is closed.
func (a *worker) StartTrigger() {
	for tr := range a.trigqueue {
		trgr, ok := a.triggers[tr.id]
		if !ok {
			glog.Errorf("Non existing trigger called: %s", tr.id)
			continue
		}
		if err := trgr.Execute(tr.objects); err != nil {
			glog.Errorf("Error execute trigger: %s", err)
		}
	}
}

// StopTrigger will stop the scaling loop.
func (a *worker) StopTrigger() {
	close(a.trigqueue)
}

// appendTrigger will append given object with given trigger ids to the given
// list of triggr objects, and will return the appended result. The result is
// a normalized list trigger id's with the corresponding objects that were
// scaled.
func (a *worker) appendTrigger(list []*triggr, obj *scanner.Object, ids []string) []*triggr {
	for _, id := range ids {
		newid := true
		for _, tr := range list {
			if tr.id == id {
				tr.objects = append(tr.objects, obj)
				newid = false
				break
			}
		}
		if newid {
			list = append(list, &triggr{id, []*scanner.Object{obj}})
		}
	}
	return list
}

// queueTriggers will enqueue the collected triggers as specified in the
// prodived list of trigger id's. Each trigger will be enqueued just once.
func (a *worker) queueTriggers(list []*triggr) {
	glog.V(6).Infof("trigger to be queueed: %#v", list)
	for _, tr := range list {
		glog.V(6).Infof("trigger added to queue: %#v", tr)
		a.trigqueue <- *tr
	}
}
