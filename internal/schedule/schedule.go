package schedule

import (
	"fmt"
	"time"
)

type Schedule struct {
	dayOfWeek map[time.Weekday]bool
	hour      int
	min       int
	org       string
	settings  map[string]string
}

// New will return a Schedule object for given schedule description.
func New(text string) (*Schedule, error) {
	s := &Schedule{
		dayOfWeek: map[time.Weekday]bool{},
		settings:  map[string]string{},
		org:       text,
	}
	if err := s.parse(text); err != nil {
		return nil, err
	}
	return s, nil
}

// GetNextTrigger will return the time the next trigger that occurs after
// given time (now) should occur according to this schedule.
func (s *Schedule) GetNextTrigger(now time.Time) (time.Time, error) {
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

// AsString will return the Schedule struct in human readable form.
func (s *Schedule) String() string {
	return s.org
}
