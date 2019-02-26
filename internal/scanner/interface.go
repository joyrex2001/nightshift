package scanner

import (
	"github.com/joyrex2001/nightswitch/internal/schedule"
)

type Scanner interface {
	GetSchedule() ([]schedule.Schedule, error)
}
