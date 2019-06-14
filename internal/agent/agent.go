package agent

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
	"github.com/joyrex2001/nightshift/internal/trigger"
)

// Agent is the public interface that is implemented by the agent.
type Agent interface {
	AddScanner(scanner.Scanner)
	AddTrigger(string, trigger.Trigger)
	SetResyncInterval(time.Duration)
	GetObjects() map[string]*scanner.Object
	GetScanners() []scanner.Scanner
	GetTriggers() map[string]trigger.Trigger
	UpdateSchedule()
	Start()
	Stop()
}

type worker struct {
	interval  time.Duration
	m         sync.Mutex
	done      chan bool
	scanners  []scanner.Scanner
	triggers  map[string]trigger.Trigger
	trigqueue chan triggr
	watchers  []watch
	objects   map[string]*objectspq
	now       time.Time
	past      time.Time
}

var instance *worker
var once sync.Once

// New will instantiate a new Agent object.
func New() Agent {
	once.Do(func() {
		instance = &worker{
			objects:   map[string]*objectspq{},
			interval:  15 * time.Minute,
			watchers:  []watch{},
			done:      make(chan bool),
			past:      time.Now().Add(-60 * time.Minute),
			scanners:  []scanner.Scanner{},
			triggers:  map[string]trigger.Trigger{},
			trigqueue: make(chan triggr, 500),
		}
	})
	return instance
}

// SetResyncInterval will set the agent resync interval to make sure that
// missing watch events are restored.
func (a *worker) SetResyncInterval(interval time.Duration) {
	a.interval = interval
}

// AddScanner will add a scanner to the agent.
func (a *worker) AddScanner(scnr scanner.Scanner) {
	a.m.Lock()
	defer a.m.Unlock()
	a.scanners = append(a.scanners, scnr)
}

// GetScanners will return the configured scanners.
func (a *worker) GetScanners() []scanner.Scanner {
	// disabled; AddScanner is only done during initialization...
	// a.m.Lock()
	// defer a.m.Unlock()
	return a.scanners
}

// AddTrigger will add a trigger to the agent.
func (a *worker) AddTrigger(id string, trgr trigger.Trigger) {
	a.m.Lock()
	defer a.m.Unlock()
	a.triggers[id] = trgr
}

// GetScanners will return the configured scanners.
func (a *worker) GetTriggers() map[string]trigger.Trigger {
	// disabled; AddTrigger is only done during initialization...
	// a.m.Lock()
	// defer a.m.Unlock()
	return a.triggers
}

// Start will start the agent.
func (a *worker) Start() {
	glog.Info("Starting agent...")
	a.UpdateSchedule()
	go a.StartWatch()
	go a.StartScale()
	go a.StartTrigger()
}

// Stop will stop the agent.
func (a *worker) Stop() {
	a.StopWatch()
	a.StopScale()
	a.StopTrigger()
}
