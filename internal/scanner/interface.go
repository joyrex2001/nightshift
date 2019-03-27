package scanner

import ()

type Scanner interface {
	GetObjects() ([]Object, error)
	SetConfig(Config)
	GetConfig() Config
	Scale(Object, int) error
}

type Scaler func(replicas int) error
