package trigger

import (
	"testing"
)

type mock struct {
	typ string
	cfg Config
	err error
}

func (m *mock) SetConfig(c Config) {
	m.cfg = c
}

func (m *mock) GetConfig() Config {
	m.cfg["type"] = m.typ
	return m.cfg
}

func (m *mock) Execute() error {
	return nil
}

func getFactory(typ string, m *mock) Factory {
	return func() (Trigger, error) {
		m.typ = typ
		m.cfg = Config{}
		return m, nil
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
			if cfg["type"] != tst.test {
				t.Errorf("failed test %d - expected %s, got: %s", i, tst.test, cfg["type"])
			}
		}
	}
}
