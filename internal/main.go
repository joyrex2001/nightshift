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
	agent := agent.New()
	if cfg := loadConfig(); cfg != nil {
		addScanners(agent, cfg)
	}
	agent.Interval = viper.GetDuration("generic.interval")
	glog.Infof("Refresh interval: %s", agent.Interval)
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

// addScanners will add configured scanners to the provided agent. The scanners
// are added in the order of priority, lowest priority is added first.
func addScanners(agent *agent.Agent, cfg *config.Config) {
	// add main config
	ns := viper.GetString("openshift.namespace")
	sel := viper.GetString("openshift.label")
	if ns != "" || sel != "" {
		scanr := scanner.NewOpenShiftScanner()
		scanr.Namespace = ns
		scanr.Label = sel
		agent.AddScanner(scanr)
	}
	// go through configured scanners
	for _, scan := range cfg.Scanner {
		def, _ := scan.Default.GetSchedule()
		// add namespace scanner
		for _, ns = range scan.Namespace {
			scanr := scanner.NewOpenShiftScanner()
			scanr.DefaultSchedule = def
			scanr.Namespace = ns
			agent.AddScanner(scanr)
		}
		// add exceptions specified in deployments
		for _, depl := range scan.Deployment {
			sched, _ := depl.GetSchedule()
			scanr := scanner.NewOpenShiftScanner()
			scanr.DefaultSchedule = def
			scanr.ForceSchedule = sched
			for _, ns = range scan.Namespace {
				scanr.Namespace = ns
				for _, sel := range depl.Selector {
					scanr.Label = sel
					agent.AddScanner(scanr)
				}
			}
		}
	}
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
