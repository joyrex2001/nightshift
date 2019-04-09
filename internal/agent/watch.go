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
	a.initWatchers()
	a.resyncScanner(quit)
	a.runWatchers()
	quit <- true
}

// StopWatch will stop watching all configured scanners.
func (a *worker) StopWatch() {
	for _, wtc := range a.watchers {
		wtc.quit <- true
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

// resyncScanner will call the UpdateSchedule method at a specified interval,
// in order to cope with missing watch events. This method will run in the
// background until the given channel will contain data.
func (a *worker) resyncScanner(quit chan bool) {
	go func() {
		for {
			tmr := time.NewTimer(a.interval)
			select {
			case <-quit:
				return
			case <-tmr.C:
				glog.V(5).Infof("Resync start...")
				a.UpdateSchedule()
				glog.V(5).Infof("Resync finisheded...")
			}
		}
	}()
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
