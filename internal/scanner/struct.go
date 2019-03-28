package scanner

import (
	"github.com/joyrex2001/nightshift/internal/schedule"
)

type Object struct {
	Namespace string               `json:"namespace"`
	UID       string               `json:"uid"`
	Name      string               `json:"name"`
	Type      ScannerType          `json:"type"`
	Schedule  []*schedule.Schedule `json:"schedule"`
	Replicas  int                  `json:"replicas"`
	Scale     Scaler               `json:"-"`
}

type Config struct {
	Namespace string               `json:"namespace"`
	Label     string               `json:"label"`
	Schedule  []*schedule.Schedule `json:"schedule"`
	Type      ScannerType          `json:"type"`
}
