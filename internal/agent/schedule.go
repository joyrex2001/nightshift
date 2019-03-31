package agent

import (
	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

func (a *worker) UpdateSchedule() {
	a.m.Lock()
	defer a.m.Unlock()
	glog.Info("Updating schedule start...")
	a.objects = map[string]*scanner.Object{}
	for _, scnr := range a.scanners {
		objs, err := scnr.GetObjects()
		if err != nil {
			glog.Errorf("Error scanning pods: %s", err)
		}
		glog.V(5).Infof("Scan result: %#v", objs)
		for _, obj := range objs {
			a.objects[obj.UID] = obj
		}
	}
	glog.V(4).Infof("Scanned objects: %v", a.objects)
	glog.V(4).Info("Updating schedule finished...")
}
