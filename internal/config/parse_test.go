package config

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func TestNew(t *testing.T) {
	_, err := New("testdata/example.yaml")
	if err != nil {
		t.Errorf("failed parsing: %s", err)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		file   string
		result *NightShift
		err    bool
	}{
		{
			file:   "testdata/nonexistingfile",
			result: &NightShift{},
			err:    true,
		},
		{
			file: "testdata/example.yaml",
			result: &NightShift{
				Scanner: []Scanner{
					Scanner{
						Namespace: []string{"development"},
						Default: Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=1",
								"Mon-Fri 18:00 replicas=0",
							},
						},
						Deployment: []Deployment{
							Deployment{
								Selector: []string{"app=shell"},
								Schedule: []string{""},
							},
						},
					},
					Scanner{
						Namespace: []string{"batch"},
						Default: Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=0",
								"Mon-Fri 18:00 replicas=1",
							},
						},
						Deployment: []Deployment{
							Deployment{
								Selector: []string{"app=shell"},
								Schedule: nil,
							},
						},
					},
				},
			},
			err: false,
		},
	}
	for i, tst := range tests {
		res, err := New(tst.file)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if !tst.err && !reflect.DeepEqual(res, tst.result) {
			t.Errorf("failed test %d - expected: %# v, got %# v", i, pretty.Formatter(tst.result), pretty.Formatter(res))
		}
	}

}
