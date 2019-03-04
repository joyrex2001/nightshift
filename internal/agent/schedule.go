package agent

import (
	"github.com/golang/glog"
	"github.com/spf13/viper"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

func (a *Agent) UpdateSchedule() {
	glog.Info("Updating schedule...")

	// TODO: itterate through scanners
	os := scanner.NewOpenShiftScanner()
	os.Namespace = viper.GetString("openshift.namespace")
	os.Label = viper.GetString("openshift.label")
	obj, err := os.GetObjects()
	if err != nil {
		glog.Errorf("Error scanning pods: %s", err)
	}
	glog.Infof("Scanned objects: %v", obj)
}
