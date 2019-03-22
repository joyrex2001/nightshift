package scanner

import (
	"fmt"
)

type ScannerType string

const (
	OpenShift ScannerType = "OpenShift"
)

// New will return a Scanner object for given ScannerType.
func New(typ ScannerType) (Scanner, error) {
	switch typ {
	case OpenShift:
		scnr := NewOpenShiftScanner()
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
