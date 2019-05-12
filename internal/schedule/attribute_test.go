package schedule

import (
	"reflect"
	"testing"
)

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

func TestGetState(t *testing.T) {
	tests := []struct {
		state State
		err   bool
		sched *Schedule
	}{
		{
			state: NoState,
			err:   false,
			sched: &Schedule{
				settings: map[string]string{},
			},
		},
		{
			state: NoState,
			err:   true,
			sched: &Schedule{
				settings: map[string]string{
					"state": "blabla",
				},
			},
		},
		{
			state: SaveState,
			err:   false,
			sched: &Schedule{
				settings: map[string]string{
					"state": "save",
				},
			},
		},
		{
			state: RestoreState,
			err:   false,
			sched: &Schedule{
				settings: map[string]string{
					"state": "restore",
				},
			},
		},
		{
			state: RestoreState,
			err:   false,
			sched: &Schedule{
				settings: map[string]string{
					"state": "rEstOre",
				},
			},
		},
		{
			state: NoState,
			err:   false,
			sched: &Schedule{
				settings: map[string]string{},
			},
		},
	}
	for i, tst := range tests {
		r, err := tst.sched.GetState()
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if r != tst.state {
			t.Errorf("failed test %d; expected %s, got %s", i, tst.state, r)
		}
	}
}

func TestGetTriggers(t *testing.T) {
	tests := []struct {
		triggers []string
		sched    *Schedule
	}{
		{
			triggers: []string{},
			sched: &Schedule{
				settings: map[string]string{},
			},
		},
		{
			triggers: []string{},
			sched: &Schedule{
				settings: map[string]string{
					"trigger": "",
				},
			},
		},
		{
			triggers: []string{"refreshdb"},
			sched: &Schedule{
				settings: map[string]string{
					"trigger": "refreshdb",
				},
			},
		},
		{
			triggers: []string{"refreshdb", "build"},
			sched: &Schedule{
				settings: map[string]string{
					"trigger": "refReshdb,build",
				},
			},
		},
		{
			triggers: []string{"refreshdb", "build"},
			sched: &Schedule{
				settings: map[string]string{
					"trigger": "refreshdb,,build",
				},
			},
		},
		{
			triggers: []string{"refreshdb", "build"},
			sched: &Schedule{
				settings: map[string]string{
					"trigger": "refreshdb,,build,",
				},
			},
		},
		{
			triggers: []string{},
			sched: &Schedule{
				settings: map[string]string{
					"trigger": ",,",
				},
			},
		},
	}
	for i, tst := range tests {
		r := tst.sched.GetTriggers()
		if !reflect.DeepEqual(r, tst.triggers) {
			t.Errorf("failed test %d; expected %#v, got %#v", i, tst.triggers, r)
		}
	}
}
