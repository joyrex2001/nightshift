package agent

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type Agent struct {
	Interval time.Duration
	m        sync.Mutex
	done     chan bool
	scanners []scanner.Scanner
	objects  map[string]scanner.Object
	now      time.Time
	past     time.Time
}

var instance *Agent
var once sync.Once

// New will instantiate a new Agent object.
func New() *Agent {
	once.Do(func() {
		instance = &Agent{
			done:     make(chan bool),
			Interval: 5 * time.Minute,
			past:     time.Now().Add(-5 * time.Minute),
			scanners: []scanner.Scanner{},
		}
	})
	return instance
}

// AddScanner will add a scanner to the agent.
func (a *Agent) AddScanner(scanner scanner.Scanner) {
	a.m.Lock()
	defer a.m.Unlock()
	a.scanners = append(a.scanners, scanner)
}

// Start will start the agent.
func (a *Agent) Start() {
	glog.Info("Starting agent...")
	go func() {
		a.loop()
	}()
}

// Stop will stop the agent.
func (a *Agent) Stop() {
	a.m.Lock()
	defer a.m.Unlock()
	a.done <- true
}

// loop will loop endlessly untile Stop has been called, calling the method
// tick at the specified Interval.
func (a *Agent) loop() {
	a.tick()
	for {
		tmr := time.NewTimer(a.Interval)
		select {
		case <-a.done:
			return
		case <-tmr.C:
			a.tick()
		}
	}
}

// tick is called at the specified Interval and will update the currentl
// configuration as specified with the given annotations, as well as Updating
// the number of replicas for deployments and statefulsets.
func (a *Agent) tick() {
	a.m.Lock()
	defer a.m.Unlock()
	a.UpdateSchedule()
	a.Scale()
}
