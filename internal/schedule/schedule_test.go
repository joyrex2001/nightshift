package schedule

import (
	"testing"
	"time"
)

func TestHasDayOfWeek(t *testing.T) {
	tests := []struct {
		day    time.Weekday
		exists bool
		sched  *Schedule
	}{
		{
			day:    0,
			exists: false,
			sched: &Schedule{
				dayOfWeek: map[time.Weekday]bool{
					1: true,
				},
			},
		},
		{
			day:    1,
			exists: true,
			sched: &Schedule{
				dayOfWeek: map[time.Weekday]bool{
					1: true,
				},
			},
		},
		{
			day:    0,
			exists: true,
			sched: &Schedule{
				dayOfWeek: map[time.Weekday]bool{
					0: true,
					1: true,
				},
			},
		},
		{
			day:    1,
			exists: true,
			sched: &Schedule{
				dayOfWeek: map[time.Weekday]bool{
					0: true,
					1: true,
				},
			},
		},
		{
			day:    2,
			exists: false,
			sched: &Schedule{
				dayOfWeek: map[time.Weekday]bool{
					0: true,
					1: true,
				},
			},
		},
	}
	for i, tst := range tests {
		if tst.sched.hasDayOfWeek(tst.day) != tst.exists {
			t.Errorf("failed test %d", i)
		}
	}
}

func TestGetReplicas(t *testing.T) {
	tests := []struct {
		replicas int
		err      bool
		sched    *Schedule
	}{
		{
			replicas: 1,
			err:      false,
			sched: &Schedule{
				settings: map[string]string{
					"replicas": "1",
				},
			},
		},
		{
			replicas: 0,
			err:      true,
			sched: &Schedule{
				settings: map[string]string{
					"replicas": "d",
				},
			},
		},
		{
			replicas: 0,
			err:      true,
			sched: &Schedule{
				settings: map[string]string{},
			},
		},
	}
	for i, tst := range tests {
		r, err := tst.sched.GetReplicas()
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if r != tst.replicas {
			t.Errorf("failed test %d; expected %d replicas, got %d", i, tst.replicas, r)
		}
	}
}
