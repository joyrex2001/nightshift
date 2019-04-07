package agent

import (
	"reflect"
	"testing"

	"github.com/joyrex2001/nightshift/internal/scanner"
)

func TestGetAddRemoveObjects(t *testing.T) {
	tests := []struct {
		add    []*scanner.Object
		remove []*scanner.Object
		result map[string]*scanner.Object
	}{
		{
			add:    []*scanner.Object{},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{},
		},
		{
			add: []*scanner.Object{
				&scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": &scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
			},
		},
		{
			add: []*scanner.Object{},
			remove: []*scanner.Object{
				&scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			result: map[string]*scanner.Object{},
		},
		{
			add: []*scanner.Object{
				&scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{
				&scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			result: map[string]*scanner.Object{},
		},
		{
			add: []*scanner.Object{
				&scanner.Object{UID: "abc", Priority: 3, Type: "myscanner"},
				&scanner.Object{UID: "abc", Priority: 1, Type: "myscanner"},
				&scanner.Object{UID: "def", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": &scanner.Object{UID: "abc", Priority: 3, Type: "myscanner"},
				"def": &scanner.Object{UID: "def", Priority: 1, Type: "myscanner"},
			},
		},
		{
			add: []*scanner.Object{
				&scanner.Object{UID: "abc", Priority: 3, Type: "myscanner"},
				&scanner.Object{UID: "abc", Priority: 3, Type: "myscanner2"},
				&scanner.Object{UID: "def", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": &scanner.Object{UID: "abc", Priority: 3, Type: "myscanner2"},
				"def": &scanner.Object{UID: "def", Priority: 1, Type: "myscanner"},
			},
		},
		{
			add: []*scanner.Object{
				&scanner.Object{UID: "abc", Priority: 1, Type: "myscanner1"},
				&scanner.Object{UID: "abc", Priority: 3, Type: "myscanner2"},
				&scanner.Object{UID: "abc", Priority: 2, Type: "myscanner3"},
				&scanner.Object{UID: "abc", Priority: 1, Type: "myscanner4"},
				&scanner.Object{UID: "abc", Priority: 3, Type: "myscanner5"},
				&scanner.Object{UID: "def", Priority: 1, Type: "myscanner6"},
				&scanner.Object{UID: "def", Priority: 2, Type: "myscanner7"},
				&scanner.Object{UID: "def", Priority: 1, Type: "myscanner8"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": &scanner.Object{UID: "abc", Priority: 3, Type: "myscanner5"},
				"def": &scanner.Object{UID: "def", Priority: 2, Type: "myscanner7"},
			},
		},
	}
	for i, tst := range tests {
		agt := &worker{}
		agt.objects = map[string]*objectspq{}
		for _, add := range tst.add {
			agt.addObject(add)
		}
		for _, rem := range tst.remove {
			agt.removeObject(rem)
		}
		objs := agt.GetObjects()
		if !reflect.DeepEqual(objs, tst.result) {
			t.Errorf("failed test %d - expected: %v, got %v", i, tst.result, objs)
		}
	}

}
