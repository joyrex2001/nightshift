package scanner

import ()

type Scanner interface {
	GetObjects() ([]Object, error)
	SetConfig(Config)
	GetConfig() Config
}

type Scaler func(replicas int) error
