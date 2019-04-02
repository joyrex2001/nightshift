package scanner

import (
	"fmt"
	"strings"
)

var modules map[string]Factory

func init() {
	modules = map[string]Factory{}
}

// RegisterModule will add the provided module, with given factory method to
// the list of available modules in order to support dependency injection, as
// well as easing up modular development for scanners.
func RegisterModule(typ string, factory Factory) {
	typ = strings.ToLower(typ)
	modules[typ] = factory
}

// New will return a Scanner object for given ScannerType.
func New(typ string) (Scanner, error) {
	typ = strings.ToLower(typ)
	factory, ok := modules[typ]
	if ok {
		scnr := factory()
		return scnr, nil
	}
	return nil, fmt.Errorf("invalid scannertype: %s", typ)
}

// NewForConfig will return a Scanner object based on the given Config object.
func NewForConfig(cfg Config) (Scanner, error) {
	scnr, err := New(cfg.Type)
	if err != nil {
		return nil, err
	}
	scnr.SetConfig(cfg)
	return scnr, nil
}
