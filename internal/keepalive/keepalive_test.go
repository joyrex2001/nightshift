package trigger

import (
	"reflect"
	"testing"
	"time"
)

func TestNewWebhookKeepAlive(t *testing.T) {
	wht, err := New()
	if err != nil {
		t.Errorf("failed test - could not instantiate webhook module; %s", err)
	}
	in := Config{
		Settings: map[string]string{
			"test1": "value1",
			"test2": "value2",
		},
	}
	wht.SetConfig(in)
	out := wht.GetConfig()
	if !reflect.DeepEqual(in, out) {
		t.Errorf("failed test - configuration not correctly set; expected %v, got %v", in, out)
	}
}

func TestTimeout(t *testing.T) {
	tests := []struct {
		cfg      Config
		duration time.Duration
		err      bool
	}{
		{
			cfg:      Config{},
			duration: 300 * time.Millisecond,
			err:      false,
		},
		{
			cfg: Config{
				Settings: map[string]string{
					"timeout": "1s",
				},
			},
			duration: 1 * time.Second,
			err:      false,
		},
		{
			cfg: Config{
				Settings: map[string]string{
					"timeout": "just wait for 1 second",
				},
			},
			duration: 0,
			err:      true,
		},
	}

	wht := &WebhookKeepAlive{}
	for i, tst := range tests {
		wht.SetConfig(tst.cfg)
		dur, err := wht.getTimeout()
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err when getTimeout: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err when getTimeout, but got none", i)
		}
		if err == nil && tst.duration != dur {
			t.Errorf("failed test %d - expected %s, but got %s", i, tst.duration, dur)
		}
		cli, err := wht.newClient()
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err when newClient: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err when newClient, but got none", i)
		}
		if err == nil && tst.duration != cli.Timeout {
			t.Errorf("failed test %d - expected %s in http.Client, but got %s", i, tst.duration, dur)
		}
	}
}
