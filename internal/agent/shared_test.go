package agent

import (
	"github.com/joyrex2001/nightshift/internal/scanner"
)

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

func (m *mockScanner) SaveState(obj *scanner.Object) (int, error) {
	m.save = true
	return 0, nil
}

func (m *mockScanner) Scale(obj *scanner.Object, r int) error {
	m.scale = r
	return nil
}

func (m *mockScanner) Watch(_stop chan bool) (chan scanner.Event, error) {
	m.out = make(chan scanner.Event)
	go func() { m.stop = <-_stop }()
	return m.out, nil
}

func getFactory(typ string, m *mockScanner) scanner.Factory {
	return func() scanner.Scanner {
		return m
	}
}
