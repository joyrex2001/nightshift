package schedule

import (
	"fmt"
	"time"
)

type Schedule struct {
	dayOfWeek map[time.Weekday]bool
	hour      int
	min       int
	settings  map[string]string
}

// New will return a Schedule object for given schedule description.
func New(text string) (*Schedule, error) {
	s := &Schedule{
		dayOfWeek: map[time.Weekday]bool{},
		settings:  map[string]string{},
	}
	if err := s.parse(text); err != nil {
		return nil, err
	}
	return s, nil
}

// GetNextTrigger will return the time the next trigger should occur according
// to this schedule.
func (s *Schedule) GetNextTrigger() (time.Time, error) {
	now := time.Now()
	next := s.getTodayTrigger()
	found := 7
	for ; now.After(next) || !s.hasDayOfWeek(next.Weekday()); found-- {
		next = next.AddDate(0, 0, 1)
	}
	if found == 0 {
		return now, fmt.Errorf("can't find next trigger, invalid schedule?")
	}
	return next, nil
}

// hasDayOfWeek checks if the given weekday is a valid configured weekday for
// this schedule.
func (s *Schedule) hasDayOfWeek(day time.Weekday) bool {
	ex, _ := s.dayOfWeek[day]
	return ex
}

// getTodayTrigger will get the trigger time if the trigger would run today.
func (s *Schedule) getTodayTrigger() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.min, 0, 0, time.Local)
}
