package scanner

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	v1beta "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1beta "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	"k8s.io/client-go/rest"
)

type StatefulSetScanner struct {
	config     Config
	kubernetes *rest.Config
}

func init() {
	RegisterModule("statefulset", NewStatefulSetScanner)
}

// NewStatefulSetScanner will instantiate a new StatefulSetScanner object.
func NewStatefulSetScanner() Scanner {
	kubernetes, err := getKubernetes()
	if err != nil {
		glog.Warningf("failed instantiating k8s client: %s", err)
	}
	return &StatefulSetScanner{
		kubernetes: kubernetes,
	}
}

// SetConfig will set the generic configuration for this scanner.
func (s *StatefulSetScanner) SetConfig(cfg Config) {
	s.config = cfg
}

// SetConfig will set the generic configuration for this scanner.
func (s *StatefulSetScanner) GetConfig() Config {
	return s.config
}

// GetObjects will return a populated list of Objects containing the relavant
// resources with their schedule info.
func (s *StatefulSetScanner) GetObjects() ([]*Object, error) {
	rcs, err := s.getStatefulSets()
	if err != nil {
		return nil, err
	}
	return s.getObjects(rcs)
}

// Scale will scale a given object to given amount of replicas.
func (s *StatefulSetScanner) Scale(obj *Object, replicas int) error {
	glog.Infof("Scaling %s/%s to %d replicas", obj.Namespace, obj.Name, replicas)
	ss, err := s.getStatefulSet(obj)
	if err != nil {
		return fmt.Errorf("GetScale failed with: %s", err)
	}
	repl := int32(replicas)
	ss.Spec.Replicas = &repl
	apps, _ := appsv1beta.NewForConfig(s.kubernetes)
	_, err = apps.StatefulSets(obj.Namespace).Update(ss)
	return nil
}

// SaveState will save the current number of replicas as an annotation on the
// statefulset config.
func (s *StatefulSetScanner) SaveState(obj *Object) error {
	ss, err := s.getStatefulSet(obj)
	if err != nil {
		return err
	}
	repl := ss.Spec.Replicas
	if ss.ObjectMeta.Annotations == nil {
		ss.ObjectMeta.Annotations = map[string]string{}
	}
	ss.ObjectMeta.Annotations[SaveStateAnnotation] = strconv.Itoa(int(*repl))
	obj.State = &State{Replicas: int(*repl)}
	apps, _ := appsv1beta.NewForConfig(s.kubernetes)
	_, err = apps.StatefulSets(obj.Namespace).Update(ss)
	return err
}

// getStatefulSet will return the statefulset for given object.
func (s *StatefulSetScanner) getStatefulSet(obj *Object) (*v1beta.StatefulSet, error) {
	if s.kubernetes == nil {
		return nil, fmt.Errorf("unable to connect to kubernetes")
	}
	apps, err := appsv1beta.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.StatefulSets(obj.Namespace).Get(obj.Name, metav1.GetOptions{})
}

// getStatefulSets will return all statefulsets in the namespace that
// match the label selector.
func (s *StatefulSetScanner) getStatefulSets() (*v1beta.StatefulSetList, error) {
	if s.kubernetes == nil {
		return nil, fmt.Errorf("unable to connect to kubernetes")
	}
	apps, err := appsv1beta.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.StatefulSets(s.config.Namespace).List(metav1.ListOptions{
		LabelSelector: s.config.Label,
	})
}

// getObjects will itterate through the list of deployment configs and populate
// a list of objects containing the schedule configuration (if any).
func (s *StatefulSetScanner) getObjects(rcs *v1beta.StatefulSetList) ([]*Object, error) {
	objs := []*Object{}
	for _, rc := range rcs.Items {
		sched, err := getSchedule(s.config.Schedule, rc.ObjectMeta.Annotations)
		if err != nil {
			glog.Errorf("error parsing schedule annotation for %s (%s); %s", rc.ObjectMeta.UID, rc.ObjectMeta.Name, err)
		}
		state, err := getState(rc.ObjectMeta.Annotations)
		if err != nil {
			glog.Errorf("error parsing state annotation for %s (%s); %s", rc.ObjectMeta.UID, rc.ObjectMeta.Name, err)
		}
		if sched != nil {
			objs = append(objs, &Object{
				Name:      rc.ObjectMeta.Name,
				Namespace: s.config.Namespace,
				UID:       string(rc.ObjectMeta.UID),
				Priority:  s.config.Priority,
				Type:      "statefulset",
				Schedule:  sched,
				State:     state,
				Replicas:  int(*rc.Spec.Replicas),
			})
		}
	}
	return objs, nil
}

// Watch will return a channel on which Event objects will be published that
// describe change events in the cluster.
func (s *StatefulSetScanner) Watch() (chan Event, error) {
	return make(chan Event), nil
}
