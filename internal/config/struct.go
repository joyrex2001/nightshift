package config

type NightShift struct {
	Scanner []Scanner `yaml:"scanner"`
}

type Scanner struct {
	Namespace  []string     `yaml:"namespace"`
	Default    Default      `yaml:"default"`
	Deployment []Deployment `yaml:"deployment"`
}

type Default struct {
	Schedule []string `yaml:"schedule"`
}

type Deployment struct {
	Selector []string `yaml:"selector"`
	Schedule []string `yaml:"schedule"`
}
