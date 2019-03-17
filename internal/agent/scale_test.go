package agent

import (
	"fmt"
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
		obj := scanner.Object{}
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
