package scanner

import (
	"github.com/joyrex2001/nightshift/internal/schedule"
)

// Config describes the configuration of a scanner. It includes ScannerType
// to allow to be used by the factory NewForConfig method.
type Config struct {
	Namespace string               `json:"namespace"`
	Label     string               `json:"label"`
	Schedule  []*schedule.Schedule `json:"schedule"`
	Type      string               `json:"type"`
}

// Object is an object found by the scanner.
type Object struct {
	Namespace string               `json:"namespace"`
	UID       string               `json:"uid"`
	Name      string               `json:"name"`
	Type      string               `json:"type"`
	Schedule  []*schedule.Schedule `json:"schedule"`
	State     *State               `json:"state"`
	Replicas  int                  `json:"replicas"`
	scanner   Scanner
}

// State defines a state of the object.
type State struct {
	Replicas int `json:"replicas"`
}
