package scanner

import (
	"errors"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/joyrex2001/nightshift/internal/schedule"
)

type mock struct {
	typ      string
	cfg      Config
	state    *Object
	scale    *Object
	replicas int
	err      error
}

func (m *mock) SetConfig(c Config) {
	m.cfg = c
}

func (m *mock) GetConfig() Config {
	m.cfg.Type = m.typ
	return m.cfg
}
func (m *mock) GetObjects() ([]*Object, error) {
	return nil, nil
}

func (m *mock) SaveState(obj *Object) (int, error) {
	m.state = obj
	return 0, nil
}

func (m *mock) Scale(obj *Object, r int) error {
	m.scale = obj
	m.replicas = r
	return m.err
}

func (m *mock) Watch(_stop chan bool) (chan Event, error) {
	return make(chan Event), nil
}

func getFactory(typ string, m *mock) Factory {
	return func() (Scanner, error) {
		m.typ = typ
		return m, nil
	}
}

func TestNewForConfig(t *testing.T) {
	m1 := &mock{}
	m2 := &mock{}
	RegisterModule("mock1", getFactory("mock1", m1))
	RegisterModule("mock2", getFactory("mock2", m2))
	tests := []struct {
		cfg  Config
		test string
		err  bool
	}{
		{Config{Type: "mock1"}, "mock1", false},
		{Config{Type: "mock2"}, "mock2", false},
		{Config{Type: "foobar"}, "foobar", true},
	}
	for i, tst := range tests {
		m, err := NewForConfig(tst.cfg)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if err == nil {
			cfg := m.GetConfig()
			if cfg.Type != tst.test {
				t.Errorf("failed test %d - expected %s, got: %s", i, tst.test, cfg.Type)
			}
		}
	}
}

func TestNew(t *testing.T) {
	m1 := &mock{}
	m2 := &mock{}
	RegisterModule("mock1", getFactory("mock1", m1))
	RegisterModule("mock2", getFactory("mock2", m2))
	tests := []struct {
		mock string
		test string
		err  bool
	}{
		{"mock1", "mock1", false},
		{"mock2", "mock2", false},
		{"foobar", "foobar", true},
	}
	for i, tst := range tests {
		m, err := New(tst.mock)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if err == nil {
			cfg := m.GetConfig()
			if cfg.Type != tst.test {
				t.Errorf("failed test %d - expected %s, got: %s", i, tst.test, cfg.Type)
			}
		}
	}
}

func TestScale(t *testing.T) {
	state := &mock{}
	RegisterModule("mock", getFactory("mock", state))
	_, err := New("mock")
	if err != nil {
		t.Errorf("failed test - unexpected err: %s", err)
	} else {
		for i := 0; i < 10; i++ {
			obj := &Object{Type: "mock"}
			state.err = nil
			err := obj.Scale(i)
			if err != nil {
				t.Errorf("failed test scaling to %d - unexpected err: %s", i, err)
			}
			if state.replicas != i {
				t.Errorf("failed test - expected %d replicas, got %d", i, state.replicas)
			}
			state.err = errors.New("some error")
			if err := obj.Scale(i); err == nil {
				t.Errorf("failed test scaling to %d - expcted an error, but got none", i)
			}
		}
	}
}

func TestSaveState(t *testing.T) {
	state := &mock{}
	RegisterModule("mock", getFactory("mock", state))
	_, err := New("mock")
	if err != nil {
		t.Errorf("failed test - unexpected err: %s", err)
	} else {
		for i := 0; i < 10; i++ {
			obj := &Object{Type: "mock", Replicas: i}
			err := obj.SaveState()
			if err != nil {
				t.Errorf("failed test %d - save state unexpected err: %s", i, err)
			}
			if state.state != obj {
				t.Error("failed test - object for save state differs")
			}
		}
	}
}

func TestNewObjectForScanner(t *testing.T) {
	scnr := &mock{typ: "mock"}
	sched := []*schedule.Schedule{{}, {}}
	cfg := Config{
		Namespace: "abc",
		Priority:  303,
		Schedule:  sched,
	}
	scnr.SetConfig(cfg)
	obj := NewObjectForScanner(scnr)
	if len(obj.Schedule) != len(sched) {
		t.Error("failed test - schedule differs")
	}
	if obj.Namespace != "abc" {
		t.Errorf("failed test - expected Namespace 'abc', got: %s", obj.Namespace)
	}
	if obj.Priority != 303 {
		t.Errorf("failed test - expected Priority '303', got: %d", obj.Priority)
	}
}

func TestUpdateWithMeta(t *testing.T) {
	obj := &Object{}
	meta := metav1.ObjectMeta{UID: "abc", Name: "something"}
	obj.updateWithMeta(meta)
	if obj.UID != "abc" {
		t.Errorf("failed test - expected UID 'abc', got: %s", obj.UID)
	}
	if obj.Name != "something" {
		t.Errorf("failed test - expected UID 'something', got: %s", obj.Name)
	}
}

func TestCopy(t *testing.T) {
	sched1, _ := schedule.New("Mon-Fri 10:00 replicas=2")
	sched2, _ := schedule.New("Thu 10:00 state=save replicas=0")
	tests := []*Object{
		{UID: "123", Name: "Something"},
		{UID: "123", Name: "Something", ScannerId: "somescanner"},
		{UID: "123", Name: "Something", State: &State{Replicas: 1}},
		{UID: "123", Name: "Something", Schedule: []*schedule.Schedule{sched1, sched2}},
	}
	for i, obj := range tests {
		new := obj.Copy()
		if new == obj {
			t.Errorf("failed test %d - objects are identical", i)
		}
		obj.UID = "changed"
		if new.UID == obj.UID {
			t.Errorf("failed test %d - change to UID is copied to new object as well", i)
		}
		if new.ScannerId != obj.ScannerId {
			t.Errorf("failed test %d - ScannerId is not copied to new object", i)
		}
		obj.ScannerId = "changed"
		if new.ScannerId == obj.ScannerId {
			t.Errorf("failed test %d - change to ScannerId is copied to new object as well", i)
		}
		if new.State != nil && new.State == obj.State {
			t.Errorf("failed test %d - object State attribute is identical (%p,%p)", i, new.State, obj.State)
		}
		if len(obj.Schedule) != len(new.Schedule) {
			t.Errorf("failed test %d - failed copying schedule length is not identical", i)
		}
		for j, sched := range obj.Schedule {
			if sched == new.Schedule[j] {
				t.Errorf("failed test %d - object Schedule attribute is identical (%p,%p)", i, sched, new.Schedule[j])
			}
			if !reflect.DeepEqual(sched, new.Schedule[j]) {
				t.Errorf("failed test %d - failed copying schedule, objects are not identical (%v vs %v)", i, sched, new.Schedule[j])
			}
		}
	}
}
