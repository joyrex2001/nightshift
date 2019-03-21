package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/joyrex2001/nightshift/internal/schedule"
)

// New will instantiate a config object for given config file. It will return
// an error if the config file is invalid, or does not exist.
func New(file string) (*Config, error) {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	m, err := loadConfig(y)
	if err != nil {
		return nil, err
	}
	if err = m.processSchedule(); err != nil {
		return nil, err
	}
	return m, nil
}

// loadConfig will load the given []byte of yaml data to a Config object.
func loadConfig(y []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(y, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// processSchedule will itterate through the config and process all schedule
// strings and cache these. It will return an error if one or more schedules
// are invalid.
func (c *Config) processSchedule() error {
	for _, scan := range c.Scanner {
		if _, err := scan.Default.GetSchedule(); err != nil {
			return err
		}
		for _, depl := range scan.Deployment {
			if _, err := depl.GetSchedule(); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetSchedule will parse the schedule strings and return an array of schedule
// objects, or an error if the schedule strings are invalid.
func (d *Default) GetSchedule() ([]*schedule.Schedule, error) {
	var err error
	if d == nil {
		return nil, nil
	}
	if d.parsed {
		return d.schedule, nil
	}
	d.parsed = true
	d.schedule, err = parseSchedule(d.Schedule)
	return d.schedule, err
}

// GetSchedule will parse the schedule strings and return an array of schedule
// objects, or an error if the schedule strings are invalid.
func (d *Deployment) GetSchedule() ([]*schedule.Schedule, error) {
	var err error
	if d.parsed {
		return d.schedule, nil
	}
	d.parsed = true
	d.schedule, err = parseSchedule(d.Schedule)
	return d.schedule, err
}

// parseSchedule will parse the schedule strings and return an array of schedule
// objects, or an error if the schedule strings are invalid.
func parseSchedule(raw []string) ([]*schedule.Schedule, error) {
	obj := []*schedule.Schedule{}
	for _, sched := range raw {
		if sched == "" {
			continue
		}
		s, err := schedule.New(sched)
		if err != nil {
			return nil, err
		}
		obj = append(obj, s)
	}
	return obj, nil
}
