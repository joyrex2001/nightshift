package internal

import (
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/joyrex2001/nightshift/internal/agent"
	"github.com/joyrex2001/nightshift/internal/config"
	"github.com/joyrex2001/nightshift/internal/webui"
)

// Main is the main entry point of this service and will start the party and
// rock the boat.
func Main(cmd *cobra.Command, args []string) {
	startAgent()
	startWebUI()
	forever()
}

// startAgent will start the agent that will monitor and scale the openshift
// resources according to the schedules.
func startAgent() {
	agent := agent.New()
	cfg := loadConfig()
	_ = cfg
	agent.Start()
}

// loadConfig will load the nightshift configuration from the configfile.
func loadConfig() *config.Config {
	if viper.ConfigFileUsed() != "" {
		cfg, err := config.New(viper.ConfigFileUsed())
		if err != nil {
			glog.Errorf("Error parsing config: %s", err)
			return nil
		}
		return cfg
	}
	return nil
}

// startWebUI will start the management webserver.
func startWebUI() {
	enabled := viper.GetBool("web.enable")
	if enabled {
		webui := webui.New()
		webui.Addr = viper.GetString("web.listen-addr")
		webui.Cert = viper.GetString("web.cert-file")
		webui.Key = viper.GetString("web.key-file")
		webui.TLS = viper.GetBool("web.enable-tls")
		webui.Start()
	}
}

// forever will wait foreva, foreva eva...
func forever() {
	for {
		time.Sleep(time.Second)
	}
}
