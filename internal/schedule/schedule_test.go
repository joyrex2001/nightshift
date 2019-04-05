package schedule

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		schedule string
		err      bool
	}{
		{
			schedule: `nonuttin`,
			err:      true,
		},
		{
			schedule: `Mon-sun 10:00 replicas=1`,
			err:      false,
		},
	}
	for i, tst := range tests {
		s, err := New(tst.schedule)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if tst.err && s != nil {
			t.Errorf("failed test %d - expected nil object", i)
		}
		if !tst.err && s == nil {
			t.Errorf("failed test %d - expected object", i)
		}
	}
}

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
		{
			day:    2,
			exists: false,
			sched: &Schedule{
				dayOfWeek: map[time.Weekday]bool{},
			},
		},
	}
	for i, tst := range tests {
		if tst.sched.hasDayOfWeek(tst.day) != tst.exists {
			t.Errorf("failed test %d", i)
		}
	}
}

func TestGetTodayTrigger(t *testing.T) {
	tests := []struct {
		timezone string
		now      time.Time
		sched    *Schedule
		trigger  time.Time
	}{
		{
			timezone: "Europe/Amsterdam",
			now:      time.Date(2019, 1, 1, 23, 0, 0, 0, time.UTC),
			sched:    &Schedule{hour: 18, min: 0},
			trigger:  time.Date(2019, 1, 2, 17, 00, 0, 0, time.UTC),
		},
		{
			timezone: "Europe/Amsterdam",
			now:      time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			sched:    &Schedule{hour: 18, min: 0},
			trigger:  time.Date(2019, 1, 1, 17, 00, 0, 0, time.UTC),
		},
	}
	for i, tst := range tests {
		SetTimeZone(tst.timezone)
		trig := tst.sched.getTodayTrigger(tst.now)
		if !trig.Equal(tst.trigger) {
			t.Errorf("failed test %d - expected time equal to %s, but got %s", i, tst.trigger, trig)
		}
	}
}

func TestGetNextTrigger(t *testing.T) {
	tests := []struct {
		timezone string
		now      time.Time
		sched    *Schedule
		trigger  time.Time
		err      bool
	}{
		{
			timezone: "UTC",
			now:      time.Date(2019, 3, 4, 8, 0, 0, 0, time.UTC), // monday
			sched: &Schedule{
				hour: 18,
				min:  0,
				dayOfWeek: map[time.Weekday]bool{
					0: true, // sunday
					1: true,
				}},
			trigger: time.Date(2019, 3, 4, 18, 0, 0, 0, time.UTC), // monday
			err:     false,
		},
		{
			timezone: "UTC",
			now:      time.Date(2019, 3, 4, 19, 0, 0, 0, time.UTC), // monday
			sched: &Schedule{
				hour: 18,
				min:  0,
				dayOfWeek: map[time.Weekday]bool{
					0: true, // sunday
					1: true,
				}},
			trigger: time.Date(2019, 3, 10, 18, 0, 0, 0, time.UTC), // monday
			err:     false,
		},
		{
			timezone: "UTC",
			now:      time.Date(2019, 3, 4, 19, 0, 0, 0, time.UTC), // monday
			sched: &Schedule{
				hour:      18,
				min:       0,
				dayOfWeek: map[time.Weekday]bool{}},
			err: true,
		},
	}
	for i, tst := range tests {
		SetTimeZone(tst.timezone)
		trig, err := tst.sched.GetNextTrigger(tst.now)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if !tst.err && !trig.Equal(tst.trigger) {
			t.Errorf("failed test %d - expected time equal to %s, but got %s", i, tst.trigger, trig)
		}
	}
}
