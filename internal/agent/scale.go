package agent

import (
	"github.com/golang/glog"
)

func (a *Agent) Scale() {
	glog.Info("Scaling resources...")
	// TODO: itterate through current schedule and update replicas accordingly
}
