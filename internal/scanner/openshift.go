package scanner

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	v1 "github.com/openshift/api/apps/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type OpenShiftScanner struct {
	config     Config
	kubernetes *rest.Config
}

const (
	OpenShift string = "OpenShift"
)

func init() {
	RegisterModule(OpenShift, NewOpenShiftScanner)
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
	cfg.Type = OpenShift
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
func (s *OpenShiftScanner) SaveState(obj *Object) error {
	dc, err := s.getDeploymentConfig(obj)
	if err != nil {
		return err
	}
	repl := dc.Spec.Replicas
	if dc.ObjectMeta.Annotations == nil {
		dc.ObjectMeta.Annotations = map[string]string{}
	}
	dc.ObjectMeta.Annotations[SaveStateAnnotation] = strconv.Itoa(int(repl))
	obj.State = &State{Replicas: int(repl)}
	apps, _ := appsv1.NewForConfig(s.kubernetes)
	_, err = apps.DeploymentConfigs(obj.Namespace).Update(dc)
	return err
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
				Type:      OpenShift,
				Schedule:  sched,
				State:     state,
				Replicas:  int(rc.Spec.Replicas),
			})
		}
	}
	return objs, nil
}
