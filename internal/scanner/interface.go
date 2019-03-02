package scanner

import (
	"github.com/joyrex2001/nightswitch/internal/schedule"
)

type Scanner interface {
	GetObjects() ([]Object, error)
}

type Object struct {
	UID      string
	Name     string
	Type     ObjectType
	Schedule []*schedule.Schedule
}

type ObjectType int

const (
	REPLICASET ObjectType = iota
	DEPLOYMENTCONFIG
)
