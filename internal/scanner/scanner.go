package scanner

import (
	"fmt"
	"strings"

	"github.com/joyrex2001/nightshift/internal/schedule"
)

// Scanner is the public interface of a scanner object.
type Scanner interface {
	SetConfig(Config)
	GetConfig() Config
	GetObjects() ([]*Object, error)
	SaveState(*Object) error
	Scale(*Object, int) error
}

// Factory is the factory method for a scanner implementation module.
type Factory func() Scanner

// Config describes the configuration of a scanner. It includes ScannerType
// to allow to be used by the factory NewForConfig method.
type Config struct {
	Namespace string               `json:"namespace"`
	Label     string               `json:"label"`
	Schedule  []*schedule.Schedule `json:"schedule"`
	Type      string               `json:"type"`
	Priority  int                  `json:"priority"`
}

// Object is an object found by the scanner.
type Object struct {
	Namespace string               `json:"namespace"`
	UID       string               `json:"uid"`
	Name      string               `json:"name"`
	Type      string               `json:"type"`
	Schedule  []*schedule.Schedule `json:"schedule"`
	State     *State               `json:"state"`
	Replicas  int                  `json:"replicas"`
	Priority  int                  `json:"priority"`
	scanner   Scanner
}

// State defines a state of the object.
type State struct {
	Replicas int `json:"replicas"`
}

var modules map[string]Factory

// RegisterModule will add the provided module, with given factory method to
// the list of available modules in order to support dependency injection, as
// well as easing up modular development for scanners.
func RegisterModule(typ string, factory Factory) {
	if modules == nil {
		modules = map[string]Factory{}
	}
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

// getScanner will lazy load the appropriate scanner object for this resource.
func (obj *Object) getScanner() (Scanner, error) {
	var err error
	if obj.scanner == nil {
		obj.scanner, err = New(obj.Type)
		if err != nil {
			return nil, err
		}
	}
	return obj.scanner, nil
}

// Scale will scale the Object to the given amount of replicas.
func (obj *Object) Scale(replicas int) error {
	scanner, err := obj.getScanner()
	if err != nil {
		return err
	}
	if err := scanner.Scale(obj, replicas); err != nil {
		return err
	}
	obj.Replicas = replicas
	return nil
}

// SaveState will save the current number of replicas.
func (obj *Object) SaveState() error {
	scanner, err := obj.getScanner()
	if err != nil {
		return err
	}
	return scanner.SaveState(obj)
}
