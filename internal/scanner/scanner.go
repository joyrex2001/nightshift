package scanner

import (
	"fmt"
	"strings"

	"github.com/joyrex2001/nightshift/internal/schedule"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Scanner is the public interface of a scanner object.
type Scanner interface {
	SetConfig(Config)
	GetConfig() Config
	GetObjects() ([]*Object, error)
	SaveState(*Object) (int, error)
	Scale(*Object, int) error
	Watch(chan bool) (chan Event, error)
}

// Factory is the factory method for a scanner implementation module.
type Factory func() (Scanner, error)

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

// Event is the structure that is send by the watch method over the channel.
type Event struct {
	Object *Object
	Type   string
}

const (
	// EventAdd is used to indicate a resource was added in a Event
	EventAdd string = "add"
	// EventRemove is used to indicate a resource was removed in a Event
	EventRemove string = "remove"
	// EventUpdate is used to indicate a resource was updated in a Event
	EventUpdate string = "update"
)

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
		return factory()
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

// NewObjectForScanner will return a new Object instance populated with the
// scanner details.
func NewObjectForScanner(scnr Scanner) *Object {
	cfg := scnr.GetConfig()
	return &Object{
		Namespace: cfg.Namespace,
		Priority:  cfg.Priority,
		Type:      cfg.Type,
		Schedule:  cfg.Schedule,
		scanner:   scnr,
	}
}

// Copy will return a fresh copy of the Object object.
func (obj *Object) Copy() *Object {
	new := &Object{}
	*new = *obj
	if new.State != nil {
		new.State = &State{}
		*(new.State) = *(obj.State)
	}
	new.Schedule = []*schedule.Schedule{}
	for _, sched := range obj.Schedule {
		new.Schedule = append(new.Schedule, sched.Copy())
	}
	return new
}

// updateWithMeta will update the Object instance with the provided kubernetes
// ObjectMeta data, and will process the supported annotations that
func (obj *Object) updateWithMeta(meta metav1.ObjectMeta) error {
	var err error
	obj.Name = meta.Name
	obj.UID = string(meta.UID)
	obj.Schedule, err = getSchedule(obj.Schedule, meta.Annotations)
	if err != nil {
		return fmt.Errorf("error parsing schedule annotation for %s (%s); %s", meta.UID, meta.Name, err)
	}
	obj.State, err = getState(meta.Annotations)
	if err != nil {
		return fmt.Errorf("error parsing state annotation for %s (%s); %s", meta.UID, meta.Name, err)
	}
	return nil
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
	repl, err := scanner.SaveState(obj)
	if err == nil {
		obj.State = &State{Replicas: repl}
	}
	return err
}
