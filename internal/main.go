package internal

import (
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/joyrex2001/nightshift/internal/agent"
	"github.com/joyrex2001/nightshift/internal/config"
	"github.com/joyrex2001/nightshift/internal/scanner"
	"github.com/joyrex2001/nightshift/internal/schedule"
	"github.com/joyrex2001/nightshift/internal/webui"
)

// Main is the main entry point of this service and will start the party and
// rock the boat.
func Main(cmd *cobra.Command, args []string) {
	// generic initialization
	tz := viper.GetString("generic.timezone")
	if err := schedule.SetTimeZone(tz); err != nil {
		glog.Errorf("Invalid timezone specified: %s", err)
	} else {
		glog.Infof("Using timezone: %s", tz)
	}
	// start subsystems
	startAgent()
	startWebUI()
	forever()
}

// startAgent will start the agent that will monitor and scale the openshift
// resources according to the schedules.
func startAgent() {
	agt := agent.New()
	if cfg := loadConfig(); cfg != nil {
		addScanners(agt, cfg)
	}
	interval := viper.GetDuration("generic.interval")
	agt.SetInterval(interval)
	glog.Infof("UpdateSchedule refresh interval: %s", interval)
	agt.Start()
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

// addScanners will add configured scanners to the provided agent. The scanners
// are added in the order of priority, lowest priority is added first.
func addScanners(agent agent.Agent, cfg *config.Config) {
	// go through configured scanners
	for _, scan := range cfg.Scanner {
		glog.V(5).Infof("Adding scanner: %v", scan)
		def, _ := scan.Default.GetSchedule()
		// add namespace scanner
		for _, ns := range scan.Namespace {
			addScanner(agent, scanner.Config{
				Type:      scan.Type,
				Namespace: ns,
				Schedule:  def,
			})
		}
		// add exceptions specified in deployments
		for _, depl := range scan.Deployment {
			sched, _ := depl.GetSchedule()
			for _, ns := range scan.Namespace {
				for _, sel := range depl.Selector {
					addScanner(agent, scanner.Config{
						Type:      scan.Type,
						Namespace: ns,
						Schedule:  sched,
						Label:     sel,
					})
				}
			}
		}
	}
	return
}

// addScanner will add a scanner specified with the scanner.Config object to
// the given agent.
func addScanner(agent agent.Agent, cfg scanner.Config) {
	scanr, err := scanner.NewForConfig(cfg)
	if err != nil {
		glog.Errorf("Error adding scanners: %s", err)
		return
	}
	agent.AddScanner(scanr)
	return
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
