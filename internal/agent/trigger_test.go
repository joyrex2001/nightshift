package agent

import (
	"reflect"
	"testing"
	"time"

	"github.com/joyrex2001/nightshift/internal/trigger"
)

func TestHandleTriggers(t *testing.T) {
	agent := &worker{}
	agent.triggers = map[string]trigger.Trigger{}
	agent.trigqueue = make(chan string)

	mock1 := &mockTrigger{}
	mock2 := &mockTrigger{}
	mock3 := &mockTrigger{}
	trigger.RegisterModule("trigger1", getTriggerFactory("trigger1", mock1))
	trigger.RegisterModule("trigger2", getTriggerFactory("trigger2", mock2))
	trigger.RegisterModule("trigger3", getTriggerFactory("trigger3", mock3))
	agent.AddTrigger("trigger1", mock1)
	agent.AddTrigger("trigger2", mock2)
	agent.AddTrigger("trigger3", mock3)

	trgrs := []string{"trigger1", "trigger1", "trigger2", "trigger1", "trigger1"}
	stopped := false
	go func() {
		agent.StartTrigger()
		stopped = true
	}()

	for _, trgr := range trgrs {
		agent.queueTrigger(trgr)
	}

	time.Sleep(time.Second)

	if mock1.exc != 4 {
		t.Errorf("invalid number of calls to trigger 1; expected 1, got %d", mock1.exc)
	}
	if mock2.exc != 1 {
		t.Errorf("invalid number of calls to trigger 2; expected 1, got %d", mock1.exc)
	}
	if mock3.exc != 0 {
		t.Errorf("invalid number of calls to trigger 3; expected 0, got %d", mock1.exc)
	}
	agent.StopTrigger()
	time.Sleep(time.Second)

	if !stopped {
		t.Errorf("StopTrigger did not stop the trigger")
	}
}

func TestQueueTriggers(t *testing.T) {
	agent := &worker{}
	agent.trigqueue = make(chan string)
	agent.triggers = map[string]trigger.Trigger{
		"trigger1": &mockTrigger{},
		"trigger2": &mockTrigger{},
	}
	res := []string{}
	trgrs := []string{"trigger1", "trigger1", "trigger2", "trigger3", "trigger1", "trigger1"}
	go agent.queueTriggers(trgrs)
	go func() {
		for trgr := range agent.trigqueue {
			res = append(res, trgr)
		}
	}()
	time.Sleep(time.Second)
	close(agent.trigqueue)

	exp := []string{"trigger1", "trigger2"}
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("failed queueTriggers - expected %s, got %s", exp, res)
	}
}
