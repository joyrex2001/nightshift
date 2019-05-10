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
	Id      string   `yaml:"id"`
	Job     *Job     `yaml:"job"`
	Webhook *Webhook `yaml:"webhook"`
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

// Job is reflection of the yaml configuration file's section "job".
type Job struct {
	Name string `yaml:"name"`
}

// Job is reflection of the yaml configuration file's section "job".
type Webhook struct {
	Url string `yaml:"url"`
}
