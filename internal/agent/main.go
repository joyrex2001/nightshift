package agent

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type Agent interface {
	AddScanner(scanner.Scanner)
	SetResyncInterval(time.Duration)
	GetObjects() map[string]*scanner.Object
	GetScanners() []scanner.Scanner
	UpdateSchedule()
	Start()
	Stop()
}

type worker struct {
	interval time.Duration
	m        sync.Mutex
	done     chan bool
	scanners []scanner.Scanner
	watchers []watch
	objects  map[string]*objectspq
	now      time.Time
	past     time.Time
}

var instance *worker
var once sync.Once

// New will instantiate a new Agent object.
func New() Agent {
	once.Do(func() {
		instance = &worker{
			objects:  map[string]*objectspq{},
			interval: 15 * time.Minute,
			watchers: []watch{},
			done:     make(chan bool),
			past:     time.Now().Add(-60 * time.Minute),
			scanners: []scanner.Scanner{},
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
func (a *worker) AddScanner(scanner scanner.Scanner) {
	a.m.Lock()
	defer a.m.Unlock()
	a.scanners = append(a.scanners, scanner)
}

// GetScanners will return the configured scanners.
func (a *worker) GetScanners() []scanner.Scanner {
	// disabled; AddScanner is only done during initialization...
	// a.m.Lock()
	// defer a.m.Unlock()
	return a.scanners
}

// Start will start the agent.
func (a *worker) Start() {
	glog.Info("Starting agent...")
	a.UpdateSchedule()
	go a.StartWatch()
	go a.StartScale()
}

// Stop will stop the agent.
func (a *worker) Stop() {
	a.StopWatch()
	a.StopScale()
}
