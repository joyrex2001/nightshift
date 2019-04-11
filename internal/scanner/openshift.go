package scanner

import (
	"fmt"

	"github.com/golang/glog"
	v1 "github.com/openshift/api/apps/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type OpenShiftScanner struct {
	config     Config
	kubernetes *rest.Config
}

func init() {
	RegisterModule("openshift", NewOpenShiftScanner)
}

// NewOpenShiftScanner will instantiate a new OpenShiftScanner object.
func NewOpenShiftScanner() Scanner {
	kubernetes, err := getKubernetes()
	if err != nil {
		glog.Warningf("failed instantiating k8s client: %s", err)
	}
	return &OpenShiftScanner{
		kubernetes: kubernetes,
	}
}

// SetConfig will set the generic configuration for this scanner.
func (s *OpenShiftScanner) SetConfig(cfg Config) {
	s.config = cfg
}

// SetConfig will set the generic configuration for this scanner.
func (s *OpenShiftScanner) GetConfig() Config {
	return s.config
}

// GetObjects will return a populated list of Objects containing the relavant
// resources with their schedule info.
func (s *OpenShiftScanner) GetObjects() ([]*Object, error) {
	rcs, err := s.getDeploymentConfigs()
	if err != nil {
		return nil, err
	}
	return s.getObjects(rcs)
}

// Scale will scale a given object to given amount of replicas.
func (s *OpenShiftScanner) Scale(obj *Object, replicas int) error {
	glog.Infof("Scaling %s/%s to %d replicas", obj.Namespace, obj.Name, replicas)
	if s.kubernetes == nil {
		return fmt.Errorf("unable to connect to kubernetes")
	}
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return err
	}
	scale, err := apps.DeploymentConfigs(obj.Namespace).GetScale(obj.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("GetScale failed with: %s", err)
	}
	scale.Spec.Replicas = int32(replicas)
	_, err = apps.DeploymentConfigs(obj.Namespace).UpdateScale(obj.Name, scale)
	return err
}

// SaveState will save the current number of replicas as an annotation on the
// deployment config.
func (s *OpenShiftScanner) SaveState(obj *Object) (int, error) {
	dc, err := s.getDeploymentConfig(obj)
	if err != nil {
		return 0, err
	}
	repl := int(dc.Spec.Replicas)
	dc.ObjectMeta = updateState(dc.ObjectMeta, repl)
	apps, _ := appsv1.NewForConfig(s.kubernetes)
	_, err = apps.DeploymentConfigs(obj.Namespace).Update(dc)
	return repl, err
}

// getDeploymentConfig will return an DeploymentConfig object.
func (s *OpenShiftScanner) getDeploymentConfig(obj *Object) (*v1.DeploymentConfig, error) {
	if s.kubernetes == nil {
		return nil, fmt.Errorf("unable to connect to kubernetes")
	}
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.DeploymentConfigs(obj.Namespace).Get(obj.Name, metav1.GetOptions{})
}

// getDeploymentConfigs will return all deploymentconfigs in the namespace that
// match the label selector.
func (s *OpenShiftScanner) getDeploymentConfigs() (*v1.DeploymentConfigList, error) {
	if s.kubernetes == nil {
		return nil, fmt.Errorf("unable to connect to kubernetes")
	}
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.DeploymentConfigs(s.config.Namespace).List(metav1.ListOptions{
		LabelSelector: s.config.Label,
	})
}

// getObjects will itterate through the list of deployment configs and populate
// a list of objects containing the schedule configuration (if any).
func (s *OpenShiftScanner) getObjects(rcs *v1.DeploymentConfigList) ([]*Object, error) {
	objs := []*Object{}
	for _, rc := range rcs.Items {
		if obj := s.getObject(&rc); obj.Schedule != nil {
			objs = append(objs, obj)
		}
	}
	return objs, nil
}

// Watch will return a channel on which Event objects will be published that
// describe change events in the cluster.
func (s *OpenShiftScanner) Watch(_stop chan bool) (chan Event, error) {
	if s.kubernetes == nil {
		return nil, fmt.Errorf("unable to connect to kubernetes")
	}
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	watcher, err := apps.DeploymentConfigs(s.config.Namespace).Watch(metav1.ListOptions{
		LabelSelector: s.config.Label,
	})

	out := make(chan Event)
	go func() {
		inch := watcher.ResultChan()
		for {
			select {
			case evt := <-inch:
				glog.V(5).Infof("Received event: %v", evt)
				s.handleEvent(out, evt)
			case <-_stop:
				return
			}
		}
	}()

	return out, nil
}

// handleEvent will take a watch event and transform it to a scanner watch
// event, and publish it to the out channel.
func (s *OpenShiftScanner) handleEvent(out chan Event, evt watch.Event) {
	if evt.Type == watch.Error {
		glog.Errorf("Error watching: %v", evt)
		return
	}

	dc, ok := evt.Object.(*v1.DeploymentConfig)
	if !ok {
		glog.Errorf("Unexpected type; %v", dc)
		return
	}

	obj := s.getObject(dc)
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

// getObject will convert a deploymentconfig object to a scanner.Object.
func (s *OpenShiftScanner) getObject(rc *v1.DeploymentConfig) *Object {
	obj := NewObjectForScanner(s)
	if err := obj.updateForMeta(rc.ObjectMeta); err != nil {
		glog.Error(err)
	}
	obj.Replicas = int(rc.Spec.Replicas)
	return obj
}
