package scanner

import (
	"fmt"

	"github.com/golang/glog"
	v1beta "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	appsv1beta "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	"k8s.io/client-go/rest"
)

// StatefulSetScanner is the object that implements scanning of OpenShift/k8s
// statefulsets.
type StatefulSetScanner struct {
	config     Config
	kubernetes *rest.Config
}

func init() {
	RegisterModule("statefulset", NewStatefulSetScanner)
}

// NewStatefulSetScanner will instantiate a new StatefulSetScanner object.
func NewStatefulSetScanner() (Scanner, error) {
	kubernetes, err := getKubernetes()
	if err != nil {
		return nil, fmt.Errorf("failed instantiating k8s client: %s", err)
	}
	return &StatefulSetScanner{
		kubernetes: kubernetes,
	}, nil
}

// SetConfig will set the generic configuration for this scanner.
func (s *StatefulSetScanner) SetConfig(cfg Config) {
	s.config = cfg
}

// GetConfig will return the config applied for this scanner.
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
func (s *StatefulSetScanner) Scale(obj *Object, state *int, replicas int) error {
	glog.Infof("Scaling %s/%s to %d replicas", obj.Namespace, obj.Name, replicas)
	ss, err := s.getStatefulSet(obj)
	if err != nil {
		return fmt.Errorf("GetScale failed with: %s", err)
	}
	repl := int32(replicas)
	ss.Spec.Replicas = &repl
	if state != nil {
		ss.ObjectMeta = updateState(ss.ObjectMeta, *state)
	}
	apps, _ := appsv1beta.NewForConfig(s.kubernetes)
	_, err = apps.StatefulSets(obj.Namespace).Update(ss)
	return err
}

// GetState will save the current number of replicas.
func (s *StatefulSetScanner) GetState(obj *Object) (int, error) {
	ss, err := s.getStatefulSet(obj)
	if err != nil {
		return 0, err
	}
	repl := int(*ss.Spec.Replicas)
	return repl, err
}

// getStatefulSet will return the statefulset for given object.
func (s *StatefulSetScanner) getStatefulSet(obj *Object) (*v1beta.StatefulSet, error) {
	apps, err := appsv1beta.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.StatefulSets(obj.Namespace).Get(obj.Name, metav1.GetOptions{})
}

// getStatefulSets will return all statefulsets in the namespace that
// match the label selector.
func (s *StatefulSetScanner) getStatefulSets() (*v1beta.StatefulSetList, error) {
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
		obj, err := s.unmarshall(&rc)
		if err != nil {
			return nil, err
		}
		if obj.Schedule != nil {
			objs = append(objs, obj)
		}
	}
	return objs, nil
}

// Watch will return a channel on which Event objects will be published that
// describe change events in the cluster.
func (s *StatefulSetScanner) Watch(_stop chan bool) (chan Event, error) {
	return watcher(_stop, s.getWatcher, s.unmarshall)
}

// getWatcher will return a watcher for DeploymentConfigs
func (s *StatefulSetScanner) getWatcher() (watch.Interface, error) {
	apps, err := appsv1beta.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.StatefulSets(s.config.Namespace).Watch(metav1.ListOptions{
		LabelSelector: s.config.Label,
	})
}

// unmarshall will convert a statefulset object to a scanner.Object.
func (s *StatefulSetScanner) unmarshall(kobj interface{}) (*Object, error) {
	m, ok := kobj.(*v1beta.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("can't unmarshall %v to Statefulset", m)
	}
	obj := NewObjectForScanner(s)
	if err := obj.updateWithMeta(m.ObjectMeta); err != nil {
		glog.Error(err)
	}
	obj.Replicas = int(*m.Spec.Replicas)
	return obj, nil
}
