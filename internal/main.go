package internal

import (
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/joyrex2001/nightshift/internal/agent"
	"github.com/joyrex2001/nightshift/internal/webui"
)

// Main is the main entry point of this service and will start the party and
// rock the boat.
func Main(cmd *cobra.Command, args []string) {
	glog.Info("Starting service...")
	agent := agent.New()
	agent.Start()
	webui := webui.New()
	webui.Start()
	forever()
}

// forever will wait foreva, foreva eva...
func forever() {
	for {
		time.Sleep(time.Second)
	}
}
