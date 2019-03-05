package config

import (
	"github.com/joyrex2001/nightshift/internal/schedule"
)

type Config struct {
	Scanner []Scanner `yaml:"scanner"`
}

type Scanner struct {
	Namespace  []string      `yaml:"namespace"`
	Default    *Default      `yaml:"default"`
	Deployment []*Deployment `yaml:"deployment"`
}

type Default struct {
	Schedule []string `yaml:"schedule"`
	schedule []*schedule.Schedule
	parsed   bool
}

type Deployment struct {
	Selector []string `yaml:"selector"`
	Schedule []string `yaml:"schedule"`
	schedule []*schedule.Schedule
	parsed   bool
}
