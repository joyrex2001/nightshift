package agent

import (
	"github.com/joyrex2001/nightshift/internal/scanner"
	"github.com/joyrex2001/nightshift/internal/trigger"
)

// mockScanner is a generic mock for scanners
type mockScanner struct {
	id    int
	scale int
	save  bool
	stop  bool
	out   chan scanner.Event
	objs  []*scanner.Object
}

func (m *mockScanner) SetConfig(c scanner.Config) {
}

func (m *mockScanner) GetConfig() scanner.Config {
	return scanner.Config{}
}

func (m *mockScanner) GetObjects() ([]*scanner.Object, error) {
	return m.objs, nil
}

func (m *mockScanner) GetState(obj *scanner.Object) (int, error) {
	m.save = true
	return 0, nil
}

func (m *mockScanner) Scale(obj *scanner.Object, state *int, r int) error {
	m.scale = r
	return nil
}

func (m *mockScanner) Watch(_stop chan bool) (chan scanner.Event, error) {
	m.out = make(chan scanner.Event)
	go func() { m.stop = <-_stop }()
	return m.out, nil
}

func getScannerFactory(typ string, m *mockScanner) scanner.Factory {
	return func() (scanner.Scanner, error) {
		return m, nil
	}
}

// mockTrigger is a generic mock for triggers
type mockTrigger struct {
	id   string
	err  error
	exc  int
	objs []*scanner.Object
	cfg  trigger.Config
}

func (m *mockTrigger) SetConfig(c trigger.Config) {
	m.cfg = c
	m.objs = []*scanner.Object{}
}

func (m *mockTrigger) GetConfig() trigger.Config {
	return m.cfg
}

func (m *mockTrigger) Execute(objs []*scanner.Object) error {
	m.exc++
	m.objs = append(m.objs, objs...)
	return m.err
}

func getTriggerFactory(typ string, m *mockTrigger) trigger.Factory {
	return func() (trigger.Trigger, error) {
		return m, nil
	}
}
