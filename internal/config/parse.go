package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func New(file string) (*Config, error) {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	m, err := loadConfig(y)
	if err != nil {
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
