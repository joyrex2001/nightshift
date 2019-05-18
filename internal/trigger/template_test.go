package trigger

import (
	"os"
	"testing"
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
				"key1": "value1",
				"key2": "value2",
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
	}

	for i, tst := range tests {
		tst.setup()
		out, err := RenderTemplate(tst.in, tst.values)
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
