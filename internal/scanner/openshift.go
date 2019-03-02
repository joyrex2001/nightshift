package scanner

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	v1 "github.com/openshift/api/apps/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/joyrex2001/nightswitch/internal/schedule"
)

type OpenShiftScanner struct {
	Namespace  string
	Label      string
	kubernetes *rest.Config
}

// New will instantiate a new scanner object for given namespace and label
// to specify which resources to scan.
func NewOpenShiftScanner(namespace, label string) Scanner {
	kubernetes, err := getKubernetes()
	if err != nil {
		glog.Warning("failed instantiating k8s client: %s", err)
	}

	return &OpenShiftScanner{
		Namespace:  namespace,
		Label:      label,
		kubernetes: kubernetes,
	}
}

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

// GetObjects will return a populated list of Objects containing the relavant
// resources with their schedule info.
func (s *OpenShiftScanner) GetObjects() ([]Object, error) {
	rcs, err := s.getDeploymentConfigs()
	if err != nil {
		return nil, err
	}
	return s.getObjects(rcs)
}

// getObjects will itterate through the list of replication controllers and
// populate a list of objects containing the schedule configuration (if any).
func (s *OpenShiftScanner) getObjects(rcs *v1.DeploymentConfigList) ([]Object, error) {
	objs := []Object{}
	for _, rc := range rcs.Items {
		ann, _ := rc.ObjectMeta.Annotations["joyrex2001.com/nightshift.schedule"]
		sched, err := s.annotationToSchedule(ann)
		if err != nil {
			glog.Errorf("error parsing schedule annotation '%s' for %s (%s); %s", ann, rc.ObjectMeta.UID, rc.ObjectMeta.Name, err)
		}
		objs = append(objs, Object{
			Name:     rc.ObjectMeta.Name,
			UID:      string(rc.ObjectMeta.UID),
			Type:     DEPLOYMENTCONFIG,
			Schedule: sched,
		})
	}
	return objs, nil
}

// getDeploymentConfigs will return all replication controllers in the
// namespace that match the label selector.
func (s *OpenShiftScanner) getDeploymentConfigs() (*v1.DeploymentConfigList, error) {
	if s.kubernetes == nil {
		return nil, fmt.Errorf("unable to connect to kubernetes")
	}

	apps, err := appsv1.NewForConfig(s.kubernetes)
	if err != nil {
		return nil, err
	}

	return apps.DeploymentConfigs(s.Namespace).List(metav1.ListOptions{
		LabelSelector: s.Label,
	})
}

// annotationToSchedule will convert the contents of the schedule annotation
// to an array of Schedule objects. It will produce an error if the provided
// annotation value is invalid.
func (s *OpenShiftScanner) annotationToSchedule(annotation string) ([]*schedule.Schedule, error) {
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
