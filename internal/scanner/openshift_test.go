package scanner

import (
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
		config := Config{
			Schedule: tst.sched,
		}
		os := &OpenShiftScanner{
			config: config,
		}
		res, err := os.getSchedule(tst.data)
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
