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
	obj   *scanner.Object
	sched *schedule.Schedule
}

// Scale will process all scanned objects and scale them accordingly.
func (a *worker) Scale() {
	a.m.Lock()
	defer a.m.Unlock()
	glog.Info("Scaling resources start...")
	a.now = time.Now()
	for _, obj := range a.objects {
		for _, e := range a.getEvents(obj) {
			glog.V(4).Infof("Scale event: %v", e)
			a.handleState(e)
			a.scale(e)
		}
	}
	a.past = a.now
	glog.V(4).Info("Scaling resources finished...")
}

// getEvents will return the events in chronological order that have to be
// done for the given object in the current tick.
func (a *worker) getEvents(obj *scanner.Object) []*event {
	var err error
	ev := []*event{}
	for _, s := range obj.Schedule {
		for next := a.past; !next.After(a.now); next = next.AddDate(0, 0, 1) {
			next, err = s.GetNextTrigger(next)
			if err != nil {
				glog.Errorf("Error processing trigger: %s", err)
				continue
			}
			if a.now.After(next) || a.now == next {
				ev = append(ev, &event{next, obj, s})
			}
		}
	}
	// order events by time
	sort.Slice(ev, func(i, j int) bool { return ev[i].at.Before(ev[j].at) })
	return ev
}

// handleState will save or restore state if this is defined in the schedule.
func (a *worker) handleState(e *event) {
	state, err := e.sched.GetState()
	if err != nil {
		glog.Errorf("Error scaling deployment: %s", err)
		return
	}
	// Save the current number of pods
	if state == schedule.SaveState {
		if err := e.obj.SaveState(); err != nil {
			glog.Errorf("Error saving state: %s", err)
			return
		}
	}
	// Restore the number of pods previously saved, and update object with the
	// State that should be applied.
	if state == schedule.RestoreState {
		if err := e.obj.LoadState(); err != nil {
			glog.Errorf("Error loading state: %s", err)
			return
		}
	}
	return
}

// scale will scale according to the event details.
func (a *worker) scale(e *event) {
	var repl int
	var err error
	if e.obj.State != nil {
		repl = e.obj.State.Replicas
	} else {
		repl, err = e.sched.GetReplicas()
	}
	if err == nil {
		err = e.obj.Scale(repl)
	}
	if err != nil {
		glog.Errorf("Error scaling deployment: %s", err)
	}
}
