package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func New(file string) (*NightShift, error) {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	m, err := LoadNightShift(y)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// LoadModel will load the given []byte of yaml data to a Model info object.
func LoadNightShift(y []byte) (*NightShift, error) {
	ns := &NightShift{}
	err := yaml.Unmarshal(y, ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}
