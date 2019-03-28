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
	Type      ScannerType          `json:"type"`
}

// Object is an object found by the scanner.
type Object struct {
	Namespace string               `json:"namespace"`
	UID       string               `json:"uid"`
	Name      string               `json:"name"`
	Type      ScannerType          `json:"type"`
	Schedule  []*schedule.Schedule `json:"schedule"`
	Replicas  int                  `json:"replicas"`
}

// Scale will scale the Object to the given amount of replicas.
func (obj Object) Scale(replicas int) error {
	scanner, err := New(obj.Type)
	if err != nil {
		return err
	}
	if err := scanner.Scale(obj, replicas); err != nil {
		return err
	}
	return nil
}
