package scanner

import (
	"github.com/joyrex2001/nightshift/internal/schedule"
)

type Scanner interface {
	GetObjects() ([]Object, error)
}

type Scaler func(replicas int) error

type Object struct {
	Namespace string
	UID       string
	Name      string
	Type      ObjectType
	Schedule  []*schedule.Schedule
	Scale     Scaler
}

type ObjectType int

const (
	DeploymentConfig ObjectType = iota
)
