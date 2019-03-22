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
	Scale     Scaler               `json:"-"`
}

type Config struct {
	Namespace       string               `json:"namespace"`
	Label           string               `json:"label"`
	ForceSchedule   []*schedule.Schedule `json:"force"`
	DefaultSchedule []*schedule.Schedule `json:"default"`
	Type            ScannerType          `json:"type"`
}
