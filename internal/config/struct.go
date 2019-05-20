package config

import (
	"github.com/joyrex2001/nightshift/internal/schedule"
)

// Config is reflection of the yaml root configuration entrypoint.
type Config struct {
	Trigger []*Trigger `yaml:"trigger"`
	Scanner []*Scanner `yaml:"scanner"`
}

// Scanner is reflection of the yaml configuration file's section "scanner".
type Scanner struct {
	Namespace  []string      `yaml:"namespace"`
	Default    *Default      `yaml:"default"`
	Deployment []*Deployment `yaml:"deployment"`
	Type       string        `yaml:"type"`
}

// Trigger is reflection of the yaml configuration file's section "trigger".
type Trigger struct {
	Id     string            `yaml:"id"`
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

// Default is reflection of the yaml configuration file's section "default".
type Default struct {
	Id       string
	Schedule []string `yaml:"schedule"`
	schedule []*schedule.Schedule
	parsed   bool
}

// Deployment is reflection of the yaml configuration file's section
// "deployment".
type Deployment struct {
	Id       string
	Selector []string `yaml:"selector"`
	Schedule []string `yaml:"schedule"`
	schedule []*schedule.Schedule
	parsed   bool
}
