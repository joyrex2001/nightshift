package scanner

import (
	"fmt"

	"github.com/golang/glog"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// DeploymentScanner is the object that implements scanning of kubernetes
// Deployments.
type DeploymentScanner struct {
	config     Config
	kubernetes *rest.Config
}

func init() {
	RegisterModule("deployment", NewDeploymentScanner)
}

// NewDeploymentScanner will instantiate a new DeploymentScanner object.
func NewDeploymentScanner() (Scanner, error) {
	kubernetes, err := getKubernetes()
	if err != nil {
		return nil, fmt.Errorf("failed instantiating k8s client: %s", err)
	}
	return &DeploymentScanner{
		kubernetes: kubernetes,
	}, nil
}

// SetConfig will set the generic configuration for this scanner.
func (s *DeploymentScanner) SetConfig(cfg Config) {
	s.config = cfg
}

// GetConfig will return the config applied for this scanner.
func (s *DeploymentScanner) GetConfig() Config {
	return s.config
}

// GetObjects will return a populated list of Objects containing the relavant
// resources with their schedule info.
func (s *DeploymentScanner) GetObjects() ([]*Object, error) {
	rcs, err := s.getDeployments()
	if err != nil {
		return nil, err
	}
	return s.getObjects(rcs)
}

// Scale will scale a given object to given amount of replicas.
func (s *DeploymentScanner) Scale(obj *Object, replicas int) error {
	glog.Infof("Scaling %s/%s to %d replicas", obj.Namespace, obj.Name, replicas)
	apps, err := kubernetes.NewForConfig(s.kubernetes)
	if err != nil {
		return err
	}
	scale, err := apps.AppsV1().Deployments(obj.Namespace).GetScale(obj.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("GetScale failed with: %s", err)
	}
	scale.Spec.Replicas = int32(replicas)
	_, err = apps.AppsV1().Deployments(obj.Namespace).UpdateScale(obj.Name, scale)
	return err
}

// SaveState will save the current number of replicas as an annotation on the
// deployment config.
func (s *DeploymentScanner) SaveState(obj *Object) (int, error) {
	dc, err := s.getDeployment(obj)
	if err != nil {
		return 0, err
	}
	repl := int(*dc.Spec.Replicas)
	dc.ObjectMeta = updateState(dc.ObjectMeta, repl)
	apps, _ := kubernetes.NewForConfig(s.kubernetes)
	_, err = apps.AppsV1().Deployments(obj.Namespace).Update(dc)
	return repl, err
}

// getDeployment will return an Deployment object.
func (s *DeploymentScanner) getDeployment(obj *Object) (*v1.Deployment, error) {
	apps, err := kubernetes.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.AppsV1().Deployments(obj.Namespace).Get(obj.Name, metav1.GetOptions{})
}

// getDeployments will return all deploymentconfigs in the namespace that
// match the label selector.
func (s *DeploymentScanner) getDeployments() (*v1.DeploymentList, error) {
	apps, err := kubernetes.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.AppsV1().Deployments(s.config.Namespace).List(metav1.ListOptions{
		LabelSelector: s.config.Label,
	})
}

// getObjects will itterate through the list of deployment configs and populate
// a list of objects containing the schedule configuration (if any).
func (s *DeploymentScanner) getObjects(rcs *v1.DeploymentList) ([]*Object, error) {
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
func (s *DeploymentScanner) Watch(_stop chan bool) (chan Event, error) {
	return watcher(_stop, s.getWatcher, s.unmarshall)
}

// getWatcher will return a watcher for Deployments
func (s *DeploymentScanner) getWatcher() (watch.Interface, error) {
	apps, err := kubernetes.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}
	return apps.AppsV1().Deployments(s.config.Namespace).Watch(metav1.ListOptions{
		LabelSelector: s.config.Label,
	})
}

// getObject will convert a deploymentconfig object to a scanner.Object.
func (s *DeploymentScanner) unmarshall(kobj interface{}) (*Object, error) {
	m, ok := kobj.(*v1.Deployment)
	if !ok {
		return nil, fmt.Errorf("can't unmarshall %v to Deployment", m)
	}
	obj := NewObjectForScanner(s)
	if err := obj.updateWithMeta(m.ObjectMeta); err != nil {
		glog.Error(err)
	}
	obj.Replicas = int(*m.Spec.Replicas)
	return obj, nil
}
