package agent

import (
	"github.com/joyrex2001/nightshift/internal/scanner"
)

type mockScanner struct {
	id    int
	scale int
	save  bool
}

func (m *mockScanner) SetConfig(c scanner.Config) {
}

func (m *mockScanner) GetConfig() scanner.Config {
	return scanner.Config{}
}

func (m *mockScanner) GetObjects() ([]*scanner.Object, error) {
	return nil, nil
}

func (m *mockScanner) SaveState(obj *scanner.Object) error {
	m.save = true
	return nil
}

func (m *mockScanner) Scale(obj *scanner.Object, r int) error {
	m.scale = r
	return nil
}

func (m *mockScanner) Watch(_stop chan bool) (chan scanner.Event, error) {
	return nil, nil
}

func getFactory(typ string, m *mockScanner) scanner.Factory {
	return func() scanner.Scanner {
		return m
	}
}
