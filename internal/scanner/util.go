package scanner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/viper"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/joyrex2001/nightshift/internal/schedule"
)

const (
	ScheduleAnnotation  string = "joyrex2001.com/nightshift.schedule"
	IgnoreAnnotation    string = "joyrex2001.com/nightshift.ignore"
	SaveStateAnnotation string = "joyrex2001.com/nightshift.savestate"
)

// getKubernetes will return a kubernetes config object.
func getKubernetes() (*rest.Config, error) {
	kubeconfig := viper.GetString("openshift.kubeconfig")
	if kubeconfig != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err == nil {
			return config, nil
		}
	}
	return rest.InClusterConfig()
}

// getState will return a State object based on the value of the State
// annotation on the deployment config. If no annotation exist, it will return
// nil.
func getState(annotations map[string]string) (*State, error) {
	repls, ok := annotations[SaveStateAnnotation]
	if !ok {
		glog.V(5).Info("no previous state available")
		return nil, nil
	}
	repl, err := strconv.Atoi(repls)
	if err != nil {
		return nil, err
	}
	return &State{Replicas: repl}, nil
}

// getSchedule will return a list of schedules, taken the annotations and
// defaults into account.
func getSchedule(cfgsched []*schedule.Schedule, annotations map[string]string) ([]*schedule.Schedule, error) {
	dis := strings.ToLower(annotations[IgnoreAnnotation])
	if dis == "true" {
		return nil, nil
	} else if dis != "false" && dis != "" {
		return nil, fmt.Errorf("invalid value '%s' for %s", dis, IgnoreAnnotation)
	}
	if ann := annotations[ScheduleAnnotation]; ann != "" {
		return annotationToSchedule(ann)
	}
	return cfgsched, nil
}

// annotationToSchedule will convert the contents of the schedule annotation
// to an array of Schedule objects. It will produce an error if the provided
// annotation value is invalid.
func annotationToSchedule(annotation string) ([]*schedule.Schedule, error) {
	sched := []*schedule.Schedule{}
	for _, ann := range strings.Split(annotation, ";") {
		if ann == "" {
			continue
		}
		s, err := schedule.New(ann)
		if err != nil {
			return nil, err
		}
		sched = append(sched, s)
	}
	return sched, nil
}
