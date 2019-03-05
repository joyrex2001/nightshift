package agent

import (
	"github.com/golang/glog"
)

func (a *Agent) UpdateSchedule() {
	glog.Info("Updating schedule...")

	for _, scnr := range a.scanners {
		obj, err := scnr.GetObjects()
		if err != nil {
			glog.Errorf("Error scanning pods: %s", err)
		}
		glog.Infof("Scanned objects: %#v", obj)
		// TODO: map objects to hashmap, override values
	}
}
