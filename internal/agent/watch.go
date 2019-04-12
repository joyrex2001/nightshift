package agent

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type watch struct {
	event chan scanner.Event
	quit  chan bool
	_quit chan bool // channel that will signal the scanner to stop watching
}

// StartWatch will start watching all configured scanners.
func (a *worker) StartWatch() {
	quit := make(chan bool)
	go a.resyncScanner(quit)
	a.initWatchers()
	a.runWatchers()
	quit <- true
}

// StopWatch will stop watching all configured scanners.
func (a *worker) StopWatch() {
	for _, wtc := range a.watchers {
		wtc.quit <- true
	}
}

// UpdateSchedule will call all scanners and get the current list of matched
// objects. This method is called periodically by the resyncScanner method to
// make sure the known state reflects the actual state of the platform, and
// makes the agent resilient against missed watch events due to e.g. network
// connectivity problems.
func (a *worker) UpdateSchedule() {
	a.InitObjects()
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

// initWatchers will initialize the watchers for all available channels.
func (a *worker) initWatchers() {
	a.watchers = []watch{}
	for _, scnr := range a.GetScanners() {
		_quit := make(chan bool)
		wtc, err := scnr.Watch(_quit)
		if err != nil {
			glog.Errorf("Error initialising watcher for scanner: %v", scnr.GetConfig())
		} else {
			a.watchers = append(a.watchers, watch{wtc, make(chan bool), _quit})
		}
	}
}

// runWatchers will run the watchers, and will block until the watchers are
// stopped by calling StopWatch().
func (a *worker) runWatchers() {
	var wg sync.WaitGroup
	for _, _wtc := range a.watchers {
		wtc := _wtc
		wg.Add(1)
		go func() {
			a.watchScanner(wtc)
			wg.Done()
		}()
	}
	wg.Wait()
}

// watchScanner will read the watch channel as provided by the scanners Watch
// method, and will update the objects according to the events received on the
// channel. This method blocks until the it receives a quit message on the
// quit channel.
func (a *worker) watchScanner(wtc watch) {
	for {
		select {
		case <-wtc.quit:
			wtc._quit <- true
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

// resyncScanner will call the UpdateSchedule method at a specified interval,
// in order to cope with missing watch events. This method will run until the
// given channel will contain data.
func (a *worker) resyncScanner(quit chan bool) {
	for {
		tmr := time.NewTimer(a.interval)
		select {
		case <-quit:
			return
		case <-tmr.C:
			glog.V(4).Infof("Resync start...")
			a.UpdateSchedule()
			glog.V(4).Infof("Resync finisheded...")
		}
	}
}
