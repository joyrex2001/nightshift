package schedule

import (
	"fmt"
	"time"
)

var timezone *time.Location

func init() {
	SetTimeZone("UTC")
}

// New will return a Schedule object for given schedule description.
func New(text string) (*Schedule, error) {
	s := &Schedule{
		dayOfWeek:   map[time.Weekday]bool{},
		settings:    map[string]string{},
		Description: text,
	}
	if err := s.parse(text); err != nil {
		return nil, err
	}
	return s, nil
}

// SetTimeZone will configure the timezone in which schedules strings
// are defined.
func SetTimeZone(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err == nil {
		timezone = loc
	}
	return err
}

// Copy will return a fresh copy of the Schedule object.
func (s *Schedule) Copy() *Schedule {
	new := &Schedule{}
	*new = *s
	return new
}

// GetNextTrigger will return the time the next trigger that occurs after
// given time (now) should occur according to this schedule.
func (s *Schedule) GetNextTrigger(now time.Time) (time.Time, error) {
	next := s.getTodayTrigger(now)
	found := 8
	for ; (now.After(next) || !s.hasDayOfWeek(next.Weekday())) && found > 0; found-- {
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
func (s *Schedule) getTodayTrigger(now time.Time) time.Time {
	now = now.In(timezone)
	return time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.min, 0, 0, timezone)
}
