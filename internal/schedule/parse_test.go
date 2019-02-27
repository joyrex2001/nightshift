package schedule

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		data  string
		err   bool
		sched *Schedule
	}{
		{
			data: `Mon 18:00 replicas=0`,
			err:  false,
			sched: &Schedule{
				hour: 18,
				min:  00,
				dayOfWeek: map[time.Weekday]bool{
					1: true,
				},
				settings: map[string]string{
					"replicas": "0",
				},
			},
		},
		{
			data: `Mon-sun 10:00 replicas=1`,
			err:  false,
			sched: &Schedule{
				hour: 10,
				min:  00,
				dayOfWeek: map[time.Weekday]bool{
					1: true,
					2: true,
					3: true,
					4: true,
					5: true,
					6: true,
					0: true,
				},
				settings: map[string]string{
					"replicas": "1",
				},
			},
		},
		{
			data: `Wed, Fri 13:30 replicas=3`,
			err:  false,
			sched: &Schedule{
				hour: 13,
				min:  30,
				dayOfWeek: map[time.Weekday]bool{
					3: true,
					5: true,
				},
				settings: map[string]string{
					"replicas": "3",
				},
			},
		},
		{
			data:  `Wed, Fri 13:66 replicas=3`,
			err:   true,
			sched: &Schedule{},
		},
		{
			data:  `Wed, Fri 34:00 replicas=3`,
			err:   true,
			sched: &Schedule{},
		},
		{
			data: `WED,  Fri 13:00 replicas=4`,
			err:  false,
			sched: &Schedule{
				hour: 13,
				min:  00,
				dayOfWeek: map[time.Weekday]bool{
					3: true,
					5: true,
				},
				settings: map[string]string{
					"replicas": "4",
				},
			},
		},
		{
			data:  `WedFri 13:00 replicas=4`,
			err:   true,
			sched: &Schedule{},
		},
		{
			data: `Wed,THu-SuN 13:00 rePLICAS=4`,
			err:  false,
			sched: &Schedule{
				hour: 13,
				min:  00,
				dayOfWeek: map[time.Weekday]bool{
					3: true,
					4: true,
					5: true,
					6: true,
					0: true,
				},
				settings: map[string]string{
					"replicas": "4",
				},
			},
		},
		{
			data: `Wed, Thu-Sun 3:10 replicas=4`,
			err:  false,
			sched: &Schedule{
				hour: 3,
				min:  10,
				dayOfWeek: map[time.Weekday]bool{
					3: true,
					4: true,
					5: true,
					6: true,
					0: true,
				},
				settings: map[string]string{
					"replicas": "4",
				},
			},
		},
		{
			data:  `Thu-Wed 13:00 replicas=4`,
			err:   true,
			sched: &Schedule{},
		},
		{
			data:  `Man 18:00 replicas=10`,
			err:   true,
			sched: &Schedule{},
		},
	}
	for i, tst := range tests {
		s := &Schedule{
			dayOfWeek: map[time.Weekday]bool{},
			settings:  map[string]string{},
		}
		err := s.parse(tst.data)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if !tst.err && !reflect.DeepEqual(s, tst.sched) {
			t.Errorf("failed test %d - expected: %v, got %v", i, tst.sched, s)
		}
	}
}

func TestTrimSpaces(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  `Mon  18:00      replicas=0`,
			out: `Mon 18:00 replicas=0`,
		},
		{
			in:  `Mon 18:00 replicas=0`,
			out: `Mon 18:00 replicas=0`,
		},
		{
			in:  `  Mon  18:00      replicas=0     `,
			out: `Mon 18:00 replicas=0`,
		},
	}
	for i, tst := range tests {
		in := trimSpaces(tst.in)
		if in != tst.out {
			t.Errorf("failed test %d - expected: %s, got %s", i, tst.out, in)
		}
	}
}
