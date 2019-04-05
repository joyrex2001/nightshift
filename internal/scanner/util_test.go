package scanner

import (
	"reflect"
	"testing"

	"github.com/joyrex2001/nightshift/internal/schedule"
)

func TestGetSchedule(t *testing.T) {
	tests := []struct {
		data  map[string]string
		sched []*schedule.Schedule
		err   bool
		count int
	}{
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0`,
			},
			err:   false,
			count: 1,
		},
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0; Mon 9:00 replicas=1;`,
			},
			err:   false,
			count: 2,
		},
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Man x:00 replicas=0; Mon 9:00 replicas=1;`,
			},
			err:   true,
			count: 0,
		},
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0; Mon 9:00 replicas=1;`,
				"joyrex2001.com/nightshift.ignore":   `something`,
			},
			err:   true,
			count: 0,
		},
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0`,
				"joyrex2001.com/nightshift.ignore":   `true`,
			},
			err:   false,
			count: 0,
		},
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0; Mon 9:00 replicas=1;`,
				"joyrex2001.com/nightshift.ignore":   `True`,
			},
			err:   false,
			count: 0,
		},
		{
			sched: []*schedule.Schedule{&schedule.Schedule{}, &schedule.Schedule{}},
			data:  map[string]string{},
			err:   false,
			count: 2,
		},
		{
			sched: []*schedule.Schedule{&schedule.Schedule{}, &schedule.Schedule{}, &schedule.Schedule{}},
			data:  map[string]string{},
			err:   false,
			count: 3,
		},
		{
			sched: []*schedule.Schedule{&schedule.Schedule{}, &schedule.Schedule{}, &schedule.Schedule{}},
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0; Mon 9:00 replicas=1;`,
				"joyrex2001.com/nightshift.ignore":   `true`,
			},
			err:   false,
			count: 0,
		},
		{
			sched: []*schedule.Schedule{&schedule.Schedule{}, &schedule.Schedule{}, &schedule.Schedule{}},
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0;`,
				"joyrex2001.com/nightshift.ignore":   `false`,
			},
			err:   false,
			count: 1,
		},
		{
			sched: []*schedule.Schedule{&schedule.Schedule{}, &schedule.Schedule{}},
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0;`,
				"joyrex2001.com/nightshift.ignore":   `false`,
			},
			err:   false,
			count: 1,
		},
	}
	for i, tst := range tests {
		res, err := getSchedule(tst.sched, tst.data)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if len(res) != tst.count {
			t.Errorf("failed test %d - invalid number of results %d, instead of %d", i, len(res), tst.count)
		}
	}
}

func TestGetState(t *testing.T) {
	tests := []struct {
		data  map[string]string
		state *State
		err   bool
	}{
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.savestate": `5`,
			},
			err:   false,
			state: &State{Replicas: 5},
		},
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.savestate": `a`,
			},
			err:   true,
			state: nil,
		},
		{
			data: map[string]string{
				"joyrex2001.com/nightshift.savestate": `a`,
				"joyrex2001.com/nightshift.ignore":    `false`,
			},
			err:   true,
			state: nil,
		},
		{
			data:  map[string]string{},
			err:   false,
			state: nil,
		},
	}
	for i, tst := range tests {
		res, err := getState(tst.data)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if !tst.err && !reflect.DeepEqual(res, tst.state) {
			t.Errorf("failed test %d - expected: %v, got %v", i, tst.state, res)
		}
	}
}
