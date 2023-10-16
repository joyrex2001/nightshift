package scanner

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	v1 "github.com/openshift/api/apps/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

// OpenShiftScanner is the object that implements scanning of OpenShift
// DeploymentConfigs.
type OpenShiftScanner struct {
	config     Config
	kubernetes *rest.Config
}

func init() {
	RegisterModule("openshift", NewOpenShiftScanner)
}

// NewOpenShiftScanner will instantiate a new OpenShiftScanner object.
func NewOpenShiftScanner() (Scanner, error) {
	kubernetes, err := getKubernetes()
	if err != nil {
		return nil, fmt.Errorf("failed instantiating k8s client: %s", err)
	}
	return &OpenShiftScanner{
		kubernetes: kubernetes,
	}, nil
}

// SetConfig will set the generic configuration for this scanner.
func (s *OpenShiftScanner) SetConfig(cfg Config) {
	s.config = cfg
}

// GetConfig will return the config applied for this scanner.
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
func (s *OpenShiftScanner) Scale(obj *Object, state *int, replicas int) error {
	glog.Infof("Scaling %s/%s to %d replicas", obj.Namespace, obj.Name, replicas)
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return err
	}
	dc, err := apps.DeploymentConfigs(obj.Namespace).Get(context.Background(), obj.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("GetScale failed with: %s", err)
	}
	if state != nil {
		dc.ObjectMeta = updateState(dc.ObjectMeta, *state)
	}
	dc.Spec.Replicas = int32(replicas)
	_, err = apps.DeploymentConfigs(obj.Namespace).Update(context.Background(), dc, metav1.UpdateOptions{})
	return err
}

// GetState will save the current number of replicas.
func (s *OpenShiftScanner) GetState(obj *Object) (int, error) {
	dc, err := s.getDeploymentConfig(obj)
	if err != nil {
		return 0, err
	}
	repl := int(dc.Spec.Replicas)
	return repl, err
}

// getDeploymentConfig will return an DeploymentConfig object.
func (s *OpenShiftScanner) getDeploymentConfig(obj *Object) (*v1.DeploymentConfig, error) {
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.DeploymentConfigs(obj.Namespace).Get(context.Background(), obj.Name, metav1.GetOptions{})
}

// getDeploymentConfigs will return all deploymentconfigs in the namespace that
// match the label selector.
func (s *OpenShiftScanner) getDeploymentConfigs() (*v1.DeploymentConfigList, error) {
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.DeploymentConfigs(s.config.Namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: s.config.Label,
	})
}

// getObjects will itterate through the list of deployment configs and populate
// a list of objects containing the schedule configuration (if any).
func (s *OpenShiftScanner) getObjects(rcs *v1.DeploymentConfigList) ([]*Object, error) {
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
func (s *OpenShiftScanner) Watch(_stop chan bool) (chan Event, error) {
	return watcher(_stop, s.getWatcher, s.unmarshall)
}

// getWatcher will return a watcher for DeploymentConfigs
func (s *OpenShiftScanner) getWatcher() (watch.Interface, error) {
	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.DeploymentConfigs(s.config.Namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: s.config.Label,
	})
}

// getObject will convert a deploymentconfig object to a scanner.Object.
func (s *OpenShiftScanner) unmarshall(kobj interface{}) (*Object, error) {
	m, ok := kobj.(*v1.DeploymentConfig)
	if !ok {
		return nil, fmt.Errorf("can't unmarshall %v to DeploymentConfig", m)
	}
	obj := NewObjectForScanner(s)
	if err := obj.updateWithMeta(m.ObjectMeta); err != nil {
		glog.Error(err)
	}
	obj.Replicas = int(m.Spec.Replicas)
	return obj, nil
}
