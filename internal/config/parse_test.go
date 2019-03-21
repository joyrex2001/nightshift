package config

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func TestNew(t *testing.T) {
	tests := []struct {
		file string
		err  bool
	}{
		{
			file: "testdata/example.yaml",
			err:  false,
		},
		{
			file: "testdata/nonexistingfile",
			err:  true,
		},
		{
			file: "testdata/invalidschedule.yaml",
			err:  true,
		},
		{
			file: "testdata/nodefault.yaml",
			err:  false,
		},
	}
	for i, tst := range tests {
		_, err := New(tst.file)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		file   string
		result *Config
		err    bool
	}{
		{
			file: "testdata/example.yaml",
			result: &Config{
				Scanner: []Scanner{
					Scanner{
						Namespace: []string{"development"},
						Default: &Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=1",
								"Mon-Fri 18:00 replicas=0",
							},
						},
						Deployment: []*Deployment{
							&Deployment{
								Selector: []string{"app=shell"},
								Schedule: []string{""},
							},
						},
					},
					Scanner{
						Namespace: []string{"batch"},
						Default: &Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=0",
								"Mon-Fri 18:00 replicas=1",
							},
						},
						Deployment: []*Deployment{
							&Deployment{
								Selector: []string{"app=shell", "app=nightshift"},
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
		y, err := ioutil.ReadFile(tst.file)
		if err != nil {
			t.Errorf("failed test %d - test configfile %s does not exist", i, err)
		}
		res, err := loadConfig(y)
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
