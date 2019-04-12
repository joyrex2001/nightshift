package agent

import (
	"reflect"
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
		&scanner.Object{UID: "abc", State: &scanner.State{Replicas: 1}},
		&scanner.Object{UID: "123", State: &scanner.State{Replicas: 2}},
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
				scanner.Event{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
			},
			result: map[string]*scanner.Object{
				"abc": &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
			},
		},
		{
			event: []scanner.Event{
				scanner.Event{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
				scanner.Event{
					Object: &scanner.Object{UID: "def", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
			},
			result: map[string]*scanner.Object{
				"abc": &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
				"def": &scanner.Object{UID: "def", Priority: 1, Type: "myscanner"},
			},
		},
		{
			event: []scanner.Event{
				scanner.Event{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
				scanner.Event{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventRemove,
				},
			},
			result: map[string]*scanner.Object{},
		},
		{
			event: []scanner.Event{
				scanner.Event{
					Object: &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
					Type:   scanner.EventAdd,
				},
				scanner.Event{
					Object: &scanner.Object{UID: "abc", Priority: 2, Type: "myscanner"},
					Type:   scanner.EventUpdate,
				},
			},
			result: map[string]*scanner.Object{
				"abc": &scanner.Object{UID: "abc", Priority: 2, Type: "myscanner"},
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
		if !reflect.DeepEqual(objs, tst.result) {
			t.Errorf("failed test %d - expected: %v, got %v", i, tst.result, objs)
		}
	}

}
