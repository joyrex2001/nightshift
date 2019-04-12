package scanner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/viper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/joyrex2001/nightshift/internal/schedule"
)

const (
	// ScheduleAnnotation is the annotation used to define schedules on resources.
	ScheduleAnnotation string = "joyrex2001.com/nightshift.schedule"
	// IgnoreAnnotation is the annotation used to define a resource should be ignored.
	IgnoreAnnotation string = "joyrex2001.com/nightshift.ignore"
	// SaveStateAnnotation is the annotation used to store the state.
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

// updateState will update a kubernetes ObjectMeta struct by either adding or
// updating the savestate annotation with the given amount of replicas. It
// will return the updated struct.
func updateState(meta metav1.ObjectMeta, repl int) metav1.ObjectMeta {
	if meta.Annotations == nil {
		meta.Annotations = map[string]string{}
	}
	meta.Annotations[SaveStateAnnotation] = strconv.Itoa(repl)
	return meta
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

// publishWatchEvent will take a watch event, and scanner object. It will
// transform it to a scanner watch event, and publish it to the out channel.
func publishWatchEvent(out chan Event, obj *Object, evt watch.Event) {
	if evt.Type == watch.Error {
		glog.Errorf("Error watching: %v", evt)
		return
	}

	if evt.Type == watch.Deleted {
		out <- Event{Object: obj, Type: EventRemove}
		return
	}

	if evt.Type == watch.Added || evt.Type == watch.Modified {
		if obj.Schedule != nil {
			out <- Event{Object: obj, Type: EventAdd}
		} else {
			out <- Event{Object: obj, Type: EventRemove}
		}
	}
}
