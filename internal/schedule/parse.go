package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parse will parse the given schedule description and store its equivalent
// attributes inside the structure.
func (s *Schedule) parse(text string) error {
	var err error

	text = trimSpaces(text)
	text = strings.Replace(text, ", ", ",", -1)
	text = strings.Replace(text, "- ", "-", -1)
	text = strings.ToLower(text)
	s.Description = text
	text = strings.Replace(text, ":", " ", -1)

	flds := strings.Split(text, " ")
	if err := s.parseWeekday(flds[0]); err != nil {
		return err
	}

	s.hour, err = strconv.Atoi(flds[1])
	if err != nil || s.hour < 0 || s.hour > 23 {
		return fmt.Errorf("invalid hour %s", flds[1])
	}

	s.min, err = strconv.Atoi(flds[2])
	if err != nil || s.min < 0 || s.min > 59 {
		return fmt.Errorf("invalid hour %s", flds[2])
	}

	for _, kv := range flds[3:] {
		kvf := strings.Split(kv, "=")
		if len(kvf) != 2 {
			return fmt.Errorf("invalid setting %s", kv)
		}
		s.settings[kvf[0]] = kvf[1]
	}

	return nil
}

// parseWeekday will parse the weekday definition of the schedule description.
func (s *Schedule) parseWeekday(text string) error {
	for _, dp := range strings.Split(text, ",") {
		oc := strings.Count(dp, "-")
		if oc > 1 {
			return fmt.Errorf("invalid day range: %s", dp)
		}
		if oc == 0 {
			wd, err := getWeekday(dp)
			if err != nil {
				return err
			}
			s.dayOfWeek[wd%7] = true
		} else {
			dpf := strings.Split(dp, "-")
			wd1, err := getWeekday(dpf[0])
			if err != nil {
				return err
			}
			wd2, err := getWeekday(dpf[1])
			if err != nil {
				return err
			}
			if wd1 >= wd2 {
				return fmt.Errorf("invalid day range: %s", dp)
			}
			for ; wd1 <= wd2; wd1++ {
				s.dayOfWeek[wd1%7] = true
			}
		}
	}
	return nil
}

// getWeekday will parse given string and return equivalent day in week.
func getWeekday(text string) (time.Weekday, error) {
	var days = map[string]time.Weekday{"sun": 7, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6}
	wd, ok := days[text]
	if !ok {
		return 0, fmt.Errorf("invalid weekday: %s", text)
	}
	return wd, nil
}

// trimSpaces will remove all leading and trailing spaces, as well as reducing
// all space sequences to just 1 space.
func trimSpaces(text string) string {
	text = strings.TrimSpace(text)
	ln := 0
	for ln != len(text) {
		ln = len(text)
		text = strings.Replace(text, "  ", " ", -1)
	}
	return text
}
