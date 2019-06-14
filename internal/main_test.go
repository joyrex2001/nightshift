package internal

import (
	"reflect"
	"testing"
	"time"

	"github.com/joyrex2001/nightshift/internal/config"
	"github.com/joyrex2001/nightshift/internal/scanner"
	"github.com/joyrex2001/nightshift/internal/trigger"
)

type scinfo struct {
	typ  string
	prio int
}

type mockAgent struct {
	trgrs []string
	scnrs []scinfo
}

func NewMockAgent() *mockAgent {
	return &mockAgent{
		trgrs: []string{},
		scnrs: []scinfo{},
	}
}

func (a *mockAgent) SetResyncInterval(t time.Duration) {}
func (a *mockAgent) UpdateSchedule()                   {}
func (a *mockAgent) Start()                            {}
func (a *mockAgent) Stop()                             {}

func (a *mockAgent) AddScanner(scnr scanner.Scanner) {
	cfg := scnr.GetConfig()
	a.scnrs = append(a.scnrs, scinfo{cfg.Type, cfg.Priority})
}

func (a *mockAgent) AddTrigger(id string, trgr trigger.Trigger) {
	a.trgrs = append(a.trgrs, id)
}

func (a *mockAgent) GetObjects() map[string]*scanner.Object {
	objs := map[string]*scanner.Object{}
	return objs
}

func (a *mockAgent) GetScanners() []scanner.Scanner {
	scnrs := []scanner.Scanner{}
	return scnrs
}

func (a *mockAgent) GetTriggers() map[string]trigger.Trigger {
	res := map[string]trigger.Trigger{}
	return res
}

type mockTrigger struct {
	id  string
	cfg trigger.Config
}

func (m *mockTrigger) SetConfig(c trigger.Config)      { m.cfg = c }
func (m *mockTrigger) GetConfig() trigger.Config       { return m.cfg }
func (m *mockTrigger) Execute([]*scanner.Object) error { return nil }

func getTriggerFactory(typ string, m *mockTrigger) trigger.Factory {
	return func() (trigger.Trigger, error) {
		return m, nil
	}
}

// mockScanner is a generic mock for scanners
type mockScanner struct {
	cfg scanner.Config
}

func (m *mockScanner) SetConfig(c scanner.Config)                        { m.cfg = c }
func (m *mockScanner) GetConfig() scanner.Config                         { return m.cfg }
func (m *mockScanner) GetObjects() ([]*scanner.Object, error)            { return []*scanner.Object{}, nil }
func (m *mockScanner) SaveState(obj *scanner.Object) (int, error)        { return 0, nil }
func (m *mockScanner) Scale(obj *scanner.Object, r int) error            { return nil }
func (m *mockScanner) Watch(_stop chan bool) (chan scanner.Event, error) { return nil, nil }

func getScannerFactory(typ string, m *mockScanner) scanner.Factory {
	return func() (scanner.Scanner, error) {
		return m, nil
	}
}

/*****************************************************************************/

func TestLoadConfig(t *testing.T) {
	res := loadConfig()
	if res != nil {
		t.Errorf("expected nil instead of config due to missing configfile (%v)", res)
	}
}

func TestAddTriggers(t *testing.T) {
	tests := []struct {
		in  *config.Config
		out []string
	}{
		{
			in:  &config.Config{},
			out: []string{},
		},
		{
			in: &config.Config{
				Trigger: []*config.Trigger{
					{
						Id:   "build",
						Type: "webhook",
					},
					{
						Id:   "dummy",
						Type: "dummyerror",
					},
				},
			},
			out: []string{"build"},
		},
	}

	for i, tst := range tests {
		agt := NewMockAgent()
		addTriggers(agt, tst.in)
		if !reflect.DeepEqual(agt.trgrs, tst.out) {
			t.Errorf("failed %d - expected %v, got %v", i, tst.out, agt.trgrs)
		}
	}
}

func TestAddAgents(t *testing.T) {
	tests := []struct {
		in  *config.Config
		out []scinfo
	}{
		{
			in:  &config.Config{},
			out: []scinfo{},
		},
		{
			in: &config.Config{
				Scanner: []*config.Scanner{
					{
						Namespace: []string{"development"},
						Default: &config.Default{
							Id: "development-default",
							Schedule: []string{
								"Mon-Fri  9:00 replicas=1",
								"Mon-Fri 18:00 replicas=0",
							},
						},
						Deployment: []*config.Deployment{
							{
								Id:       "development-shell",
								Selector: []string{"app=shell"},
								Schedule: []string{""},
							},
						},
						Type: "mockscanner",
					},
					{
						Namespace: []string{"batch"},
						Default: &config.Default{
							Id: "batch",
							Schedule: []string{
								"Mon-Fri  9:00 replicas=0",
								"Mon-Fri 18:00 replicas=1",
							},
						},
						Deployment: []*config.Deployment{
							{
								Id:       "",
								Selector: []string{"app=shell", "app=nightshift"},
								Schedule: nil,
							},
						},
						Type: "mockscanner",
					},
					{
						Namespace: []string{"batch"},
						Deployment: []*config.Deployment{
							{
								Id:       "",
								Selector: []string{"app=shell", "app=nightshift"},
								Schedule: nil,
							},
						},
						Type: "mockscanner",
					},
					{
						Namespace: []string{"batch"},
						Default:   &config.Default{},
						Type:      "errorscanner",
					},
				},
			},
			out: []scinfo{{"mockscanner", 0}, {"mockscanner", 1}, {"mockscanner", 2}, {"mockscanner", 3}, {"mockscanner", 4}, {"mockscanner", 5}, {"mockscanner", 6}, {"mockscanner", 7}},
		},
	}

	mock := &mockScanner{}
	scanner.RegisterModule("mockscanner", getScannerFactory("mockscanner", mock))

	for i, tst := range tests {
		agt := NewMockAgent()
		addScanners(agt, tst.in)
		if !reflect.DeepEqual(agt.scnrs, tst.out) {
			t.Errorf("failed %d - expected %v, got %v", i, tst.out, agt.scnrs)
		}
	}
}
