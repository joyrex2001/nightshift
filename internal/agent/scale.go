package agent

import (
	"sort"
	"time"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
	"github.com/joyrex2001/nightshift/internal/schedule"
)

type event struct {
	at    time.Time
	obj   scanner.Object
	sched *schedule.Schedule
}

// Scale will process all scanned objects and scale them accordingly.
func (a *Agent) Scale() {
	glog.Info("Scaling resources start...")
	a.now = time.Now()
	for _, obj := range a.objects {
		for _, e := range a.getEvents(obj) {
			glog.Infof("Scale event: %v", e)
		}
	}
	a.past = a.now
	glog.Info("Scaling resources finished...")
}

// getEvents will return the events in chronological order that have to be
// done for the given object in the current tick.
func (a *Agent) getEvents(obj scanner.Object) []event {
	ev := []event{}
	for _, s := range obj.Schedule {
		next, err := s.GetNextTrigger(a.past)
		if err != nil {
			glog.Errorf("Error processing trigger: %s", err)
			continue
		}
		if a.now.After(next) || a.now == next {
			ev = append(ev, event{next, obj, s})
		}
	}
	// order events by time
	sort.Slice(ev, func(i, j int) bool { return ev[j].at.Before(ev[i].at) })
	return ev
}
