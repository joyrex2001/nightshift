package agent

import (
	"sync"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type watch struct {
	event chan scanner.Event
	done  chan bool
}

// StartWatch will start watching all configured scanners.
func (a *worker) StartWatch() {
	a.watchers = []watch{}
	for _, scnr := range a.GetScanners() {
		wtc, err := scnr.Watch()
		if err != nil {
			glog.Errorf("Error initialising watcher for scanner: %v", scnr.GetConfig())
		} else {
			a.watchers = append(a.watchers, watch{wtc, make(chan bool)})
		}
	}
	var wg sync.WaitGroup
	for _, wtc := range a.watchers {
		go func() {
			wg.Add(1)
			a.watchScanner(wtc)
			wg.Done()
		}()
	}
	wg.Wait()
}

// StopWatch will stop watching all configured scanners.
func (a *worker) StopWatch() {
	for _, wtc := range a.watchers {
		wtc.done <- true
	}
}

// watchScanner will read the watch channel as provided by the scanners Watch
// method, and will update the objects according to the events received on the
// channel.
func (a *worker) watchScanner(wtc watch) {
	for {
		select {
		case <-wtc.done:
			return
		case event := <-wtc.event:
			glog.V(4).Infof("Watch event: %v", event)
			if event.Type == scanner.EventRemove {
				a.removeObject(event.Object)
			} else {
				a.addObject(event.Object)
			}
		}
	}
}

// UpdateSchedule is the periodically called schedule update method. This
// method will be obsolete once the watchers are implemented and the updates
// can be handled in real time.
func (a *worker) UpdateSchedule() {
	for _, scnr := range a.GetScanners() {
		objs, err := scnr.GetObjects()
		if err != nil {
			glog.Errorf("Error scanning pods: %s", err)
		}
		glog.V(5).Infof("Scan result: %#v", objs)
		for _, obj := range objs {
			a.addObject(obj)
		}
	}
}
