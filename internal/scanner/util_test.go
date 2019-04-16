package scanner

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

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
			sched: []*schedule.Schedule{{}, {}},
			data:  map[string]string{},
			err:   false,
			count: 2,
		},
		{
			sched: []*schedule.Schedule{{}, {}, {}},
			data:  map[string]string{},
			err:   false,
			count: 3,
		},
		{
			sched: []*schedule.Schedule{{}, {}, {}},
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0; Mon 9:00 replicas=1;`,
				"joyrex2001.com/nightshift.ignore":   `true`,
			},
			err:   false,
			count: 0,
		},
		{
			sched: []*schedule.Schedule{{}, {}, {}},
			data: map[string]string{
				"joyrex2001.com/nightshift.schedule": `Mon 18:00 replicas=0;`,
				"joyrex2001.com/nightshift.ignore":   `false`,
			},
			err:   false,
			count: 1,
		},
		{
			sched: []*schedule.Schedule{{}, {}},
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

func TestUpdateState(t *testing.T) {
	meta := metav1.ObjectMeta{}
	meta = updateState(meta, 10)
	st := meta.Annotations["joyrex2001.com/nightshift.savestate"]
	if st != "10" {
		t.Errorf("failed test - expected: 10, got %s", st)
	}
	meta = updateState(meta, 5)
	st = meta.Annotations["joyrex2001.com/nightshift.savestate"]
	if st != "5" {
		t.Errorf("failed test - expected: 5, got %s", st)
	}
}

func TestPublishWatchEvent(t *testing.T) {
	sched := []*schedule.Schedule{{}}
	tests := []struct {
		out Event
		obj *Object
		evt watch.Event
	}{
		{
			out: Event{},
			obj: &Object{UID: "123"},
			evt: watch.Event{Type: watch.Error},
		},
		{
			out: Event{Object: &Object{UID: "123"}, Type: EventRemove},
			obj: &Object{UID: "123"},
			evt: watch.Event{Type: watch.Deleted},
		},
		{
			out: Event{Object: &Object{UID: "123", Schedule: sched}, Type: EventAdd},
			obj: &Object{UID: "123", Schedule: sched},
			evt: watch.Event{Type: watch.Added},
		},
		{
			out: Event{Object: &Object{UID: "123", Schedule: sched}, Type: EventAdd},
			obj: &Object{UID: "123", Schedule: sched},
			evt: watch.Event{Type: watch.Modified},
		},
		{
			out: Event{Object: &Object{UID: "123"}, Type: EventRemove},
			obj: &Object{UID: "123"},
			evt: watch.Event{Type: watch.Added},
		},
		{
			out: Event{Object: &Object{UID: "123"}, Type: EventRemove},
			obj: &Object{UID: "123"},
			evt: watch.Event{Type: watch.Modified},
		},
	}
	for i, tst := range tests {
		in := make(chan Event, 1)
		publishWatchEvent(in, tst.obj, tst.evt)
		close(in)
		out := Event{}
		for out = range in {
		}
		if !reflect.DeepEqual(out, tst.out) {
			t.Errorf("failed test %d - expected: %v, got %v", i, tst.out, out)
		}
	}
}

func TestWatcher(t *testing.T) {
	var conns int
	var doerr error

	stop := make(chan bool)
	w := watch.NewFake()

	connect := func() (watch.Interface, error) {
		conns++
		return w, doerr
	}

	unmarsh := func(interface{}) (*Object, error) {
		return &Object{}, nil
	}

	// test error connecting
	doerr = fmt.Errorf("oops")
	_, err := watcher(stop, connect, unmarsh)
	if err == nil {
		t.Errorf("failed test watcher - expected error but got none")
	}

	// check successful connect
	doerr = nil
	out, err := watcher(stop, connect, unmarsh)
	if err != nil {
		t.Errorf("failed test watcher - unexpected error: %s", err)
	}
	if conns != 2 {
		t.Errorf("failed test watcher - expected: 2 connection attemps, got %d", conns)
	}

	// test error evt
	w.Action(watch.Error, nil)
	time.Sleep(time.Second)
	if conns != 3 {
		t.Errorf("failed test watcher - expected: 3 connection attemps, got %d", conns)
	}

	// done testing
	close(out)
	for range out {
	}
	stop <- true
}
