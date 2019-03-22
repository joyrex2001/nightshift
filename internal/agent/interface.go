package agent

import (
	"time"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

type Agent interface {
	AddScanner(scanner.Scanner)
	SetInterval(time.Duration)
	GetObjects() map[string]scanner.Object
	GetScanners() []scanner.Scanner
	Start()
	Stop()
}
