package agent

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/joyrex2001/nightshift/internal/scanner"
	"github.com/joyrex2001/nightshift/internal/trigger"
)

func TestHandleTriggers(t *testing.T) {
	agent := &worker{}
	agent.triggers = map[string]trigger.Trigger{}
	agent.trigqueue = make(chan triggr, 500)

	mock1 := &mockTrigger{}
	mock2 := &mockTrigger{err: fmt.Errorf("oops")}
	mock3 := &mockTrigger{}
	trigger.RegisterModule("trigger1", getTriggerFactory("trigger1", mock1))
	trigger.RegisterModule("trigger2", getTriggerFactory("trigger2", mock2))
	trigger.RegisterModule("trigger3", getTriggerFactory("trigger3", mock3))
	agent.AddTrigger("trigger1", mock1)
	agent.AddTrigger("trigger2", mock2)
	agent.AddTrigger("trigger3", mock3)

	obj1 := &scanner.Object{}
	obj2 := &scanner.Object{}
	obj3 := &scanner.Object{}
	obj4 := &scanner.Object{}

	agent.queueTriggers(agent.appendTrigger([]*triggr{}, obj1, []string{"trigger1"}))
	agent.queueTriggers(agent.appendTrigger([]*triggr{}, obj2, []string{"trigger2"}))
	agent.queueTriggers(agent.appendTrigger([]*triggr{}, obj3, []string{"trigger1"}))
	agent.queueTriggers(agent.appendTrigger([]*triggr{}, obj4, []string{"trigger1"}))
	agent.queueTriggers(agent.appendTrigger([]*triggr{}, obj4, []string{"trigger4"}))

	stopped := false
	go func() {
		agent.StartTrigger()
		stopped = true
	}()

	time.Sleep(time.Second)

	if mock1.exc != 3 {
		t.Errorf("invalid number of calls to trigger 1; expected 3, got %d", mock1.exc)
	}
	if !reflect.DeepEqual(mock1.objs, []*scanner.Object{obj1, obj3, obj4}) {
		t.Errorf("invalid number object for trigger 1; got %v", mock1.objs)
	}
	if mock2.exc != 1 {
		t.Errorf("invalid number of calls to trigger 2; expected 1, got %d", mock2.exc)
	}
	if !reflect.DeepEqual(mock2.objs, []*scanner.Object{obj2}) {
		t.Errorf("invalid number object for trigger 2; got %v", mock2.objs)
	}
	if mock3.exc != 0 {
		t.Errorf("invalid number of calls to trigger 3; expected 0, got %d", mock3.exc)
	}
	if len(mock3.objs) != 0 {
		t.Errorf("invalid number object for trigger 3; got %v", mock3.objs)
	}
	agent.StopTrigger()
	time.Sleep(time.Second)

	if !stopped {
		t.Errorf("StopTrigger did not stop the trigger")
	}
}

func TestAppendTriggers(t *testing.T) {
	agent := &worker{}
	agent.trigqueue = make(chan triggr, 500)
	agent.triggers = map[string]trigger.Trigger{
		"trigger1": &mockTrigger{},
		"trigger2": &mockTrigger{},
	}

	obj1 := &scanner.Object{}
	obj2 := &scanner.Object{}
	obj3 := &scanner.Object{}
	obj4 := &scanner.Object{}
	obj5 := &scanner.Object{}

	trgrs := []*triggr{}
	trgrs = agent.appendTrigger(trgrs, obj1, []string{"trigger1"})
	trgrs = agent.appendTrigger(trgrs, obj2, []string{"trigger2", "trigger3"})
	trgrs = agent.appendTrigger(trgrs, obj3, []string{"trigger1"})
	trgrs = agent.appendTrigger(trgrs, obj4, []string{"trigger2"})
	trgrs = agent.appendTrigger(trgrs, obj5, []string{"trigger1"})

	exp := []*triggr{
		{id: "trigger1", objects: []*scanner.Object{obj1, obj3, obj5}},
		{id: "trigger2", objects: []*scanner.Object{obj2, obj4}},
		{id: "trigger3", objects: []*scanner.Object{obj2}},
	}

	for i, res := range exp {
		if !reflect.DeepEqual(res, trgrs[i]) {
			t.Errorf("failed appendTrigger - expected %v, got %v", res, trgrs[i])
		}
	}
}
