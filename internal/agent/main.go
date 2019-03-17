package agent

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type Agent interface {
	AddScanner(scanner.Scanner)
	SetInterval(time.Duration)
	Start()
	Stop()
}

type worker struct {
	interval time.Duration
	m        sync.Mutex
	done     chan bool
	scanners []scanner.Scanner
	objects  map[string]scanner.Object
	now      time.Time
	past     time.Time
}

var instance *worker
var once sync.Once

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

// loop will loop endlessly untile Stop has been called, calling the method
// tick at the specified interval.
func (a *worker) loop() {
	a.tick()
	for {
		tmr := time.NewTimer(a.interval)
		select {
		case <-a.done:
			return
		case <-tmr.C:
			a.tick()
		}
	}
}

// tick is called at the specified interval and will update the currentl
// configuration as specified with the given annotations, as well as Updating
// the number of replicas for deployments and statefulsets.
func (a *worker) tick() {
	a.m.Lock()
	defer a.m.Unlock()
	a.UpdateSchedule()
	a.Scale()
}
