package agent

import (
	"testing"
	"time"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

func TestStartStopWatch(t *testing.T) {
	wrkr := &worker{}
	scnr := &mockScanner{}
	wrkr.AddScanner(scnr)
	go wrkr.StartWatch()
	time.Sleep(time.Second)
	wrkr.StopWatch()
	time.Sleep(time.Second)
	if !scnr.stop {
		t.Error("scanner did not stop...")
	}
}

func TestUpdateSchedule(t *testing.T) {
	wrkr := &worker{}
	scnr := &mockScanner{}
	wrkr.AddScanner(scnr)

	scnr.objs = []*scanner.Object{
		{UID: "abc", State: &scanner.State{Replicas: 1}},
		{UID: "123", State: &scanner.State{Replicas: 2}},
	}

	wrkr.UpdateSchedule()
	objs := wrkr.GetObjects()
	if len(objs) != len(scnr.objs) {
		t.Errorf("failed test UpdateSchedule - expected: %d objects, got %d objects", len(scnr.objs), len(objs))
	}
}

func TestWatchScanner(t *testing.T) {
	tests := []struct {
		event  []scanner.Event
		result map[string]*scanner.Object
	}{
		{
			event:  []scanner.Event{},
			result: map[string]*scanner.Object{},
		},
		{
			event: []scanner.Event{
				{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
			},
			result: map[string]*scanner.Object{
				"abc": {UID: "abc", Priority: 1, Type: "myscanner"},
			},
		},
		{
			event: []scanner.Event{
				{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
				{
					Object: &scanner.Object{UID: "def", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
			},
			result: map[string]*scanner.Object{
				"abc": {UID: "abc", Priority: 1, Type: "myscanner"},
				"def": {UID: "def", Priority: 1, Type: "myscanner"},
			},
		},
		{
			event: []scanner.Event{
				{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
				{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventRemove,
				},
			},
			result: map[string]*scanner.Object{},
		},
		{
			event: []scanner.Event{
				{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
				{
					Object: &scanner.Object{UID: "abc", Priority: 2, Type: "myscanner"},
					Type:   scanner.EventUpdate,
				},
			},
			result: map[string]*scanner.Object{
				"abc": {UID: "abc", Priority: 2, Type: "myscanner"},
			},
		},
	}

	for i, tst := range tests {
		wrkr := &worker{}
		wrkr.InitObjects()

		wtc := watch{
			event: make(chan scanner.Event),
			quit:  make(chan bool),
			_quit: make(chan bool),
		}

		go wrkr.watchScanner(wtc)
		for _, evt := range tst.event {
			wtc.event <- evt
		}
		wtc.quit <- true

		objs := wrkr.GetObjects()
		for j, obj := range objs {
			if obj == tst.result[j] {
				t.Errorf("failed test %d - expected a copy, but got identical object instance", i)
			}
			if obj.UID != tst.result[j].UID {
				t.Errorf("failed test %d - expected: %v, got %v", i, tst.result[j], obj)
			}
			if obj.Priority != tst.result[j].Priority {
				t.Errorf("failed test %d - expected: %v, got %v", i, tst.result[j], obj)
			}
			if obj.Type != tst.result[j].Type {
				t.Errorf("failed test %d - expected: %v, got %v", i, tst.result[j], obj)
			}
		}
	}

}
