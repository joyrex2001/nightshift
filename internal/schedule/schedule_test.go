package schedule

import (
	"fmt"
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

func TestGetTodayTrigger(t *testing.T) {
	SetTimeZone("Europe/Amsterdam")
	s := Schedule{hour: 18, min: 0}
	t1 := time.Date(2019, 1, 1, 23, 0, 0, 0, time.UTC)
	t2 := s.getTodayTrigger(t1)
	t3 := time.Date(2019, 1, 2, 17, 30, 0, 0, time.UTC)

	fmt.Printf("t1=%s\n", t1)
	fmt.Printf("t2=%s\n", t2)
	fmt.Printf("t3=%s\n", t3)

	fmt.Printf("t2.epoch=%d\n", t2.Unix())
	fmt.Printf("t3.epoch=%d\n", t3.Unix())

	if t3.After(t2) {
		fmt.Printf("ok=========\n")
	}

}
