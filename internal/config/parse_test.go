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
			file: "testdata/invalidschedule1.yaml",
			err:  true,
		},
		{
			file: "testdata/invalidschedule2.yaml",
			err:  true,
		},
		{
			file: "testdata/invalidyaml.yaml",
			err:  true,
		},
		{
			file: "testdata/nodefault.yaml",
			err:  false,
		},
		{
			file: "testdata/empty.yaml",
			err:  false,
		},
		{
			file: "testdata/triggers.yaml",
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
				Scanner: []*Scanner{
					{
						Namespace: []string{"development"},
						Default: &Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=1",
								"Mon-Fri 18:00 replicas=0",
							},
						},
						Deployment: []*Deployment{
							{
								Selector: []string{"app=shell"},
								Schedule: []string{""},
							},
						},
						Type: "openshift",
					},
					{
						Namespace: []string{"batch"},
						Default: &Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=0",
								"Mon-Fri 18:00 replicas=1",
							},
						},
						Deployment: []*Deployment{
							{
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

func TestParseTrigger(t *testing.T) {
	tests := []struct {
		file   string
		result *Config
		err    bool
	}{
		{
			file: "testdata/triggers.yaml",
			result: &Config{
				Trigger: []*Trigger{
					{
						Id:     "cleanup",
						Type:   "job",
						Config: map[string]string{"name": "cleanup"},
					},
					{
						Id:     "refreshdb",
						Type:   "job",
						Config: map[string]string{"name": "somejob"},
					},
					{
						Id:     "build",
						Type:   "webhook",
						Config: map[string]string{"url": "http://localhost:8080", "timeout": "1s"},
					},
					{
						Id:     "dummy",
						Type:   "dummyerror",
						Config: map[string]string{},
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
		res.processTriggers()
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if !tst.err && !reflect.DeepEqual(res, tst.result) {
			t.Errorf("failed test %d - expected: %# v, got %# v", i, pretty.Formatter(tst.result), pretty.Formatter(res))
		}
	}
}

func TestProcessDefaults(t *testing.T) {
	tests := []struct {
		in  *Config
		out *Config
	}{
		{
			in: &Config{
				Scanner: []*Scanner{
					{
						Namespace: []string{"batch"},
						Default: &Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=0",
								"Mon-Fri 18:00 replicas=1",
							},
						},
						Deployment: []*Deployment{
							{
								Selector: []string{"app=shell", "app=nightshift"},
								Schedule: nil,
							},
						},
					},
				},
			},
			out: &Config{
				Scanner: []*Scanner{
					{
						Namespace: []string{"batch"},
						Default: &Default{
							Schedule: []string{
								"Mon-Fri  9:00 replicas=0",
								"Mon-Fri 18:00 replicas=1",
							},
						},
						Deployment: []*Deployment{
							{
								Selector: []string{"app=shell", "app=nightshift"},
								Schedule: nil,
							},
						},
						Type: "openshift",
					},
				},
			},
		},
	}
	for i, tst := range tests {
		tst.in.processDefaults()
		if !reflect.DeepEqual(tst.out, tst.in) {
			t.Errorf("failed test %d - expected: %# v, got %# v", i, pretty.Formatter(tst.out), pretty.Formatter(tst.in))
		}
	}
}
