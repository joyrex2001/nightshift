package agent

import (
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
				{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": {UID: "abc", Priority: 1, Type: "myscanner"},
			},
		},
		{
			add: []*scanner.Object{},
			remove: []*scanner.Object{
				{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			result: map[string]*scanner.Object{},
		},
		{
			add: []*scanner.Object{
				{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{
				{UID: "abc", Priority: 1, Type: "myscanner"},
			},
			result: map[string]*scanner.Object{},
		},
		{
			add: []*scanner.Object{
				{UID: "abc", Priority: 3, Type: "myscanner"},
				{UID: "abc", Priority: 1, Type: "myscanner"},
				{UID: "def", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": {UID: "abc", Priority: 3, Type: "myscanner"},
				"def": {UID: "def", Priority: 1, Type: "myscanner"},
			},
		},
		{
			add: []*scanner.Object{
				{UID: "abc", Priority: 3, Type: "myscanner"},
				{UID: "abc", Priority: 3, Type: "myscanner2"},
				{UID: "def", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": {UID: "abc", Priority: 3, Type: "myscanner2"},
				"def": {UID: "def", Priority: 1, Type: "myscanner"},
			},
		},
		{
			add: []*scanner.Object{
				{UID: "abc", Priority: 3, Type: "myscanner"},
				{UID: "abc", Priority: 3, Type: "myscanner2"},
				{UID: "def", Priority: 1, Type: "myscanner"},
			},
			remove: []*scanner.Object{
				{UID: "abc", Priority: 3, Type: "myscanner2"},
			},
			result: map[string]*scanner.Object{
				"def": {UID: "def", Priority: 1, Type: "myscanner"},
			},
		},
		{
			add: []*scanner.Object{
				{UID: "abc", Priority: 1, Type: "myscanner1"},
				{UID: "abc", Priority: 3, Type: "myscanner2"},
				{UID: "abc", Priority: 2, Type: "myscanner3"},
				{UID: "abc", Priority: 1, Type: "myscanner4"},
				{UID: "abc", Priority: 3, Type: "myscanner5"},
				{UID: "def", Priority: 1, Type: "myscanner6"},
				{UID: "def", Priority: 2, Type: "myscanner7"},
				{UID: "def", Priority: 1, Type: "myscanner8"},
			},
			remove: []*scanner.Object{},
			result: map[string]*scanner.Object{
				"abc": {UID: "abc", Priority: 3, Type: "myscanner5"},
				"def": {UID: "def", Priority: 2, Type: "myscanner7"},
			},
		},
	}
	for i, tst := range tests {
		agt := &worker{}
		agt.InitObjects()
		agt.objects = map[string]*objectspq{}
		for _, add := range tst.add {
			agt.addObject(add)
		}
		for _, rem := range tst.remove {
			agt.removeObject(rem)
		}
		objs := agt.GetObjects()
		for j, obj := range objs {
			if obj == tst.result[j] {
				t.Errorf("failed test %d - expected a copy, but got identical object instance", i)
			}
			if obj.UID != tst.result[j].UID {
				t.Errorf("failed test %d - expected: %v, got %v", i, tst.result[j], obj)
			}
			if obj.Priority != tst.result[j].Priority {
				t.Errorf("failed test %d - expected: %v, got %v", i, tst.result[j], obj)
			}
			if obj.Type != tst.result[j].Type {
				t.Errorf("failed test %d - expected: %v, got %v", i, tst.result[j], obj)
			}
		}
	}

}
