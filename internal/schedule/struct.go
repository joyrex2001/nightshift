package schedule

import (
	"time"
)

// Schedule is the object that describes the schedule. It contains one
// read-only attribute "Description". To actually use it, use one of the
// public methods on this object.
type Schedule struct {
	Description string `json:"Description"`
	dayOfWeek   map[time.Weekday]bool
	hour        int
	min         int
	settings    map[string]string
}

// State describes the possible values of the 'state' attribute.
type State string

var (
	RestoreState State = "restore"
	SaveState    State = "save"
	NoState      State = ""
)
