package scanner

import (
	"github.com/joyrex2001/nightshift/internal/schedule"
)

type Scanner interface {
	GetObjects() ([]Object, error)
}

type Object struct {
	Namespace string
	UID       string
	Name      string
	Type      ObjectType
	Schedule  []*schedule.Schedule
}

type ObjectType int

const (
	DeploymentConfig ObjectType = iota
)
