package agent

import (
	"github.com/golang/glog"
	"github.com/spf13/viper"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

func (a *Agent) UpdateSchedule() {
	glog.Info("Updating schedule...")

	// TODO: itterate through scanners
	ns := viper.GetString("openshift.namespace")
	lb := viper.GetString("openshift.label")

	os := scanner.NewOpenShiftScanner(ns, lb)
	obj, err := os.GetObjects()
	if err != nil {
		glog.Errorf("error scanning pods: %s", err)
	}
	glog.Infof("Scanned objects: %#v", obj)
}
