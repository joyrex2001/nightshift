package agent

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type worker struct {
	interval time.Duration
	m        sync.Mutex
	done     chan bool
	scanners []scanner.Scanner
	objects  map[string]*objectspq
	now      time.Time
	past     time.Time
}

var instance *worker
var once sync.Once

const scaleInterval = 30 * time.Second

// New will instantiate a new Agent object.
func New() Agent {
	once.Do(func() {
		instance = &worker{
			done:     make(chan bool),
			interval: 5 * time.Minute,
			past:     time.Now().Add(-60 * time.Minute),
			scanners: []scanner.Scanner{},
		}
	})
	return instance
}

// SetInterval will set the agent refresh interval.
func (a *worker) SetInterval(interval time.Duration) {
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
	go func() {
		a.loop()
	}()
}

// Stop will stop the agent.
func (a *worker) Stop() {
	a.m.Lock()
	defer a.m.Unlock()
	a.done <- true
}

// loop will loop endlessly untile Stop has been called, calling the Scale and
// UpdateSchedule methods at a specified interval.
func (a *worker) loop() {
	// Make sure everything is updated when starting the tick loop.
	a.UpdateSchedule()
	a.Scale()

	sched := time.NewTimer(a.interval)
	scale := time.NewTimer(scaleInterval)
	for {
		select {
		case <-a.done:
			return
		case <-sched.C:
			a.UpdateSchedule()
			sched.Reset(a.interval)
		case <-scale.C:
			a.Scale()
			scale.Reset(scaleInterval)
		}
	}
}
