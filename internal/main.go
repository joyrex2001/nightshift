package internal

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/joyrex2001/nightswitch/internal/agent"
	"github.com/joyrex2001/nightswitch/internal/webui"
)

func Main(cmd *cobra.Command, args []string) {
	glog.Info("Starting service...")
	agent := agent.New()
	agent.Start()
	webui.Main()
}
