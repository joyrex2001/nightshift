package schedule

import (
	"time"
)

type Schedule struct {
	dayOfWeek map[int]bool
	hour      int
	min       int
	settings  map[string]string
}

// New will return a Schedule object for given schedule description.
func New(text string) (*Schedule, error) {
	s := &Schedule{
		dayOfWeek: map[int]bool{},
		settings:  map[string]string{},
	}
	if err := s.parse(text); err != nil {
		return nil, err
	}
	return s, nil
}

// GetNextTrigger will return the time the next trigger should occur according
// to this schedule.
func (s *Schedule) GetNextTrigger() time.Time {

	return time.Now()
}
