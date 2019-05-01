package schedule

import (
	"fmt"
	"strconv"
	"strings"
)

// GetReplicas will return the number of replicas that should be applied
// according to the schedule.
func (s *Schedule) GetReplicas() (int, error) {
	r, ok := s.settings["replicas"]
	if !ok {
		return 0, fmt.Errorf("replicas definition not found in schedule")
	}
	return strconv.Atoi(r)
}

// GetState will return the state that should be applied according to the
// schedule.
func (s *Schedule) GetState() (State, error) {
	r, ok := s.settings["state"]
	if !ok {
		return NoState, nil
	}
	st, ok := map[string]State{
		"save":    SaveState,
		"restore": RestoreState,
	}[strings.ToLower(r)]
	if !ok {
		return NoState, fmt.Errorf("invalid state provided: %s", r)
	}
	return st, nil
}

// GetTriggers will return the reference codes of the triggers that should be
// triggered.
func (s *Schedule) GetTriggers() []string {
	trgs := []string{}
	for _, trg := range strings.Split(s.settings["trigger"], ",") {
		if trg != "" {
			trgs = append(trgs, trg)
		}
	}
	return trgs
}
