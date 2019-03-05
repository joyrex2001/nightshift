package agent

import (
	"github.com/golang/glog"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

func (a *Agent) UpdateSchedule() {
	glog.Info("Updating schedule start...")
	a.objects = map[string]scanner.Object{}
	for _, scnr := range a.scanners {
		objs, err := scnr.GetObjects()
		if err != nil {
			glog.Errorf("Error scanning pods: %s", err)
		}
		glog.V(4).Infof("Scan result: %#v", objs)
		for _, obj := range objs {
			a.objects[obj.UID] = obj
		}
	}
	glog.Infof("Scanned objects: %v", a.objects)
	glog.Info("Updating schedule finished...")
}
