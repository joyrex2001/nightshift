package trigger

import (
	"os"
	"testing"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		in     string
		out    string
		setup  func()
		values Config
		err    bool
	}{
		{
			in:     "",
			out:    "",
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `environment var a contains: {{ env "a" }}`,
			out:    `environment var a contains: bc`,
			values: Config{},
			setup: func() {
				os.Setenv("a", "bc")
			},
			err: false,
		},
		{
			in:  `config key1 contains {{ .key1 }} and config key2 contains {{ .key2 }}`,
			out: `config key1 contains value1 and config key2 contains value2`,
			values: Config{
				Settings: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			setup: func() {},
			err:   false,
		},
		{
			in:  `config key1 contains {{ .key1 }} and matched namespace contains {{ .objects.foo.Namespace }} with type {{ .objects.foo.Type }}`,
			out: `config key1 contains value1 and matched namespace contains default with type deploymentconfig`,
			values: Config{
				Settings: map[string]string{
					"key1": "value1",
				},
				Objects: map[string]*scanner.Object{
					"foo": &scanner.Object{
						Namespace: "default",
						Name:      "app",
						Type:      "deploymentconfig",
					},
				},
			},
			setup: func() {},
			err:   false,
		},
		{
			in:     "malformed template with an unknown {{ function }}",
			out:    "",
			values: Config{},
			setup:  func() {},
			err:    true,
		},
		{
			in:     `10 + 20 = {{ add "10" "20" }}`,
			out:    `10 + 20 = 30`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `10 + -20 = {{ add "10" "-20" }}`,
			out:    `10 + -20 = -10`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `no numbers a + a = {{ add "a" "a" }}`,
			out:    `no numbers a + a = 0`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `epoch 0 = {{ time "rfc3339" "0" }}`,
			out:    `epoch 0 = 1970-01-01T00:00:00Z`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `epoch 0 = {{ add "0" "10" | time "ansic" }}`,
			out:    `epoch 0 = Thu Jan  1 00:00:10 1970`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `epoch 0 = {{ add "0" "20" | time "unixdate" }}`,
			out:    `epoch 0 = Thu Jan  1 00:00:20 UTC 1970`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `epoch 0 + 10 = {{ add "0" "10" | time "20060102150405" }}`,
			out:    `epoch 0 + 10 = 19700101000010`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `now test {{- now | time "" }}`,
			out:    `now test`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
		{
			in:     `invalid time = {{ time "rfc3339" "a" }}`,
			out:    `invalid time = 1970-01-01T00:00:00Z`,
			values: Config{},
			setup:  func() {},
			err:    false,
		},
	}

	os.Setenv("TZ", "UTC")
	for i, tst := range tests {
		tst.setup()
		tstTempVars := mixinObjects(tst.values.Settings, tst.values.Objects)
		out, err := RenderTemplate(tst.in, tstTempVars)
		if err != nil && !tst.err {
			t.Errorf("failed test %d - unexpected err: %s", i, err)
		}
		if err == nil && tst.err {
			t.Errorf("failed test %d - expected err, but got none", i)
		}
		if err == nil && tst.out != out {
			t.Errorf("failed test %d - expected %s, but got %s", i, tst.out, out)
		}
	}
}
