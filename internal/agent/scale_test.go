package agent

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/joyrex2001/nightshift/internal/scanner"
	"github.com/joyrex2001/nightshift/internal/schedule"
)

func TestGetEvents(t *testing.T) {
	tests := []struct {
		past   time.Time
		now    time.Time
		sched  []string
		events []time.Time
	}{
		{
			past: time.Date(2019, 3, 4, 0, 0, 0, 0, time.UTC), // monday
			now:  time.Date(2019, 3, 10, 23, 59, 0, 0, time.UTC),
			sched: []string{
				"Mon-Fri 8:00 replicas=1",
				"Mon-Fri 18:00 replicas=0",
				"Sat,Sun 14:00 replicas=1",
				"Sat 16:00 replicas=0",
				"Sun 15:00 replicas=0",
			},
			events: []time.Time{
				time.Date(2019, 3, 4, 8, 0, 0, 0, time.UTC), // monday
				time.Date(2019, 3, 4, 18, 0, 0, 0, time.UTC),
				time.Date(2019, 3, 5, 8, 0, 0, 0, time.UTC), // tuesday
				time.Date(2019, 3, 5, 18, 0, 0, 0, time.UTC),
				time.Date(2019, 3, 6, 8, 0, 0, 0, time.UTC), // wednesday
				time.Date(2019, 3, 6, 18, 0, 0, 0, time.UTC),
				time.Date(2019, 3, 7, 8, 0, 0, 0, time.UTC), // thursday
				time.Date(2019, 3, 7, 18, 0, 0, 0, time.UTC),
				time.Date(2019, 3, 8, 8, 0, 0, 0, time.UTC), // friday
				time.Date(2019, 3, 8, 18, 0, 0, 0, time.UTC),
				time.Date(2019, 3, 9, 14, 0, 0, 0, time.UTC), // saturday
				time.Date(2019, 3, 9, 16, 0, 0, 0, time.UTC),
				time.Date(2019, 3, 10, 14, 0, 0, 0, time.UTC), // sunday
				time.Date(2019, 3, 10, 15, 0, 0, 0, time.UTC),
			},
		},
		{
			past: time.Date(2019, 3, 4, 8, 0, 0, 0, time.UTC), // monday
			now:  time.Date(2019, 3, 4, 8, 1, 0, 0, time.UTC),
			sched: []string{
				"Mon-Fri 8:00 replicas=1",
				"Mon-Fri 18:00 replicas=0",
				"Sat,Sun 14:00 replicas=1",
				"Sat 16:00 replicas=0",
				"Sun 15:00 replicas=0",
			},
			events: []time.Time{
				time.Date(2019, 3, 4, 8, 0, 0, 0, time.UTC), // monday
			},
		},
		{
			past: time.Date(2019, 3, 4, 10, 0, 0, 0, time.UTC), // monday
			now:  time.Date(2019, 3, 4, 11, 0, 0, 0, time.UTC),
			sched: []string{
				"Mon-Fri 8:00 replicas=1",
			},
			events: []time.Time{},
		},
	}

	for i, tst := range tests {
		agt := &worker{}
		agt.past = tst.past
		agt.now = tst.now
		obj := &scanner.Object{}
		obj.Schedule = []*schedule.Schedule{}
		for _, s := range tst.sched {
			sc, err := schedule.New(s)
			if err != nil {
				t.Errorf("failed test %d - unexpected error parsing schedule: %s", i, err)
			} else {
				obj.Schedule = append(obj.Schedule, sc)
			}
		}
		evts := agt.getEvents(obj)
		for j, evt := range evts {
			fmt.Printf("[%02d] %s\n", j, evt.at)
		}
		if len(evts) != len(tst.events) {
			t.Errorf("failed test %d - invalid number of events, expected: %v, got %v", i, len(tst.events), len(evts))
		} else {
			for j, evt := range evts {
				if evt.at != tst.events[j] {
					t.Errorf("failed test %d.%d - invalid events, expected: %v, got %v", i, j, tst.events[j], evt.at)
				}
			}
		}
	}
}

func TestHandleStateScale(t *testing.T) {
	mock := &mockScanner{}
	scanner.RegisterModule("scanner", getScannerFactory("scanner", mock))

	tests := []struct {
		sched   string
		obj     *scanner.Object
		restore bool
		save    bool
		scale   int
	}{
		{
			sched:   "Mon-Fri 8:00 replicas=3 state=restore",
			obj:     &scanner.Object{State: &scanner.State{Replicas: 1}},
			restore: true,
			save:    false,
			scale:   1,
		},
		{
			sched:   "Mon-Fri 8:00 replicas=3 state=restore",
			obj:     &scanner.Object{},
			restore: false,
			save:    false,
			scale:   3,
		},
		{
			sched:   "Mon-Fri 8:00 replicas=2 state=save",
			obj:     &scanner.Object{},
			restore: false,
			save:    true,
			scale:   2,
		},
	}

	for i, tst := range tests {
		agent := &worker{}
		tst.obj.Type = "scanner"
		sc, _ := schedule.New(tst.sched)
		tst.obj.Schedule = []*schedule.Schedule{sc}
		mock.save = false

		evt := &event{
			obj:     tst.obj,
			sched:   sc,
			restore: false,
		}

		agent.handleState(evt)
		if evt.restore != tst.restore {
			t.Errorf("failed test %d - invalid state handling restore, expected: %v, got %v", i, tst.restore, evt.restore)
		}
		if mock.save != tst.save {
			t.Errorf("failed test %d - invalid state handling save, expected: %v, got %v", i, tst.save, mock.save)
		}

		agent.scale(evt)
		if mock.scale != tst.scale {
			t.Errorf("failed test %d - invalid scaling, expected: %d replicas, got %d replicas", i, tst.scale, mock.scale)
		}
	}

}

func TestAppendEventTriggers(t *testing.T) {
	tests := []struct {
		sched string
		init  []string
		trgrs []string
	}{
		{
			init:  []string{},
			sched: "Mon-Fri 8:00 replicas=3 state=restore trigger=trigger2,trigger3",
			trgrs: []string{"trigger2", "trigger3"},
		},
		{
			init:  []string{"trigger1"},
			sched: "Mon-Fri 8:00 replicas=3 state=restore trigger=trigger2,trigger3",
			trgrs: []string{"trigger1", "trigger2", "trigger3"},
		},
		{
			init:  []string{"trigger1", "trigger2"},
			sched: "Mon-Fri 8:00 replicas=3 state=restore trigger=trigger2,trigger3",
			trgrs: []string{"trigger1", "trigger2", "trigger2", "trigger3"},
		},
	}

	for i, tst := range tests {
		agent := &worker{}
		sc, _ := schedule.New(tst.sched)
		evt := &event{sched: sc}
		trgrs := agent.appendEventTriggers(tst.init, evt)
		if !reflect.DeepEqual(trgrs, tst.trgrs) {
			t.Errorf("failed test %d - expected %s, got %s", i, tst.trgrs, trgrs)
		}
	}
}
