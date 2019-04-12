package config

import (
	"github.com/joyrex2001/nightshift/internal/schedule"
)

// Config is reflection of the yaml root configuration entrypoint.
type Config struct {
	Scanner []*Scanner `yaml:"scanner"`
}

// Scanner is reflection of the yaml configuration file's section "scanner".
type Scanner struct {
	Namespace  []string      `yaml:"namespace"`
	Default    *Default      `yaml:"default"`
	Deployment []*Deployment `yaml:"deployment"`
	Type       string        `yaml:"type"`
}

// Default is reflection of the yaml configuration file's section "default".
type Default struct {
	Schedule []string `yaml:"schedule"`
	schedule []*schedule.Schedule
	parsed   bool
}

// Deployment is reflection of the yaml configuration file's section
// "deployment".
type Deployment struct {
	Selector []string `yaml:"selector"`
	Schedule []string `yaml:"schedule"`
	schedule []*schedule.Schedule
	parsed   bool
}
