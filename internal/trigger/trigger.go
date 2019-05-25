package trigger

import (
	"fmt"
	"strings"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

// Trigger defines the public interface of trigger modules.
type Trigger interface {
	SetConfig(Config)
	GetConfig() Config
	Execute() error
}

// Config is the configuration for this trigger, and contains a hashmap with
// generic settings. The key for each value should be lowercased always.
type Config struct {
	Id       string                     `json:"id"`
	Type     string                     `json:"type"`
	Settings map[string]string          `json:"settings"`
	Objects  map[string]*scanner.Object `json:"agentObjects"`
}

// Factory is the factory method for a trigger implementation module.
type Factory func() (Trigger, error)

var modules map[string]Factory

// RegisterModule will add the provided module, with given factory method to
// the list of available modules in order to support dependency injection, as
// well as easing up modular development for triggers.
func RegisterModule(typ string, factory Factory) {
	if modules == nil {
		modules = map[string]Factory{}
	}
	typ = strings.ToLower(typ)
	modules[typ] = factory
}

// New will return a Trigger object for given TriggerType.
func New(typ string) (Trigger, error) {
	typ = strings.ToLower(typ)
	factory, ok := modules[typ]
	if ok {
		return factory()
	}
	return nil, fmt.Errorf("invalid triggertype: %s", typ)
}
