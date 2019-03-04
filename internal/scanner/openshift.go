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

	"github.com/joyrex2001/nightshift/internal/schedule"
)

type OpenShiftScanner struct {
	Namespace       string
	Label           string
	ForceSchedule   []*schedule.Schedule
	DefaultSchedule []*schedule.Schedule
	kubernetes      *rest.Config
}

// NewOpenShiftScanner will instantiate a new OpenShiftScanner object.
func NewOpenShiftScanner() *OpenShiftScanner {
	kubernetes, err := getKubernetes()
	if err != nil {
		glog.Warning("failed instantiating k8s client: %s", err)
	}
	return &OpenShiftScanner{
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

// getObjects will itterate through the list of replication controllers and
// populate a list of objects containing the schedule configuration (if any).
func (s *OpenShiftScanner) getObjects(rcs *v1.DeploymentConfigList) ([]Object, error) {
	objs := []Object{}
	for _, rc := range rcs.Items {
		sched, err := s.getSchedule(rc.ObjectMeta.Annotations)
		if err != nil {
			glog.Errorf("error parsing schedule annotation for %s (%s); %s", rc.ObjectMeta.UID, rc.ObjectMeta.Name, err)
		}
		if sched != nil {
			objs = append(objs, Object{
				Name:      rc.ObjectMeta.Name,
				Namespace: s.Namespace,
				UID:       string(rc.ObjectMeta.UID),
				Type:      DeploymentConfig,
				Schedule:  sched,
			})
		}
	}
	return objs, nil
}

// getSchedule will return a list of schedules, taken the annotations and
// defaults into account.
func (s *OpenShiftScanner) getSchedule(annotations map[string]string) ([]*schedule.Schedule, error) {
	dis := strings.ToLower(annotations["joyrex2001.com/nightshift.ignore"])
	if dis == "true" {
		return nil, nil
	} else if dis != "false" && dis != "" {
		return nil, fmt.Errorf("invalid value '%s' for nightshift.ignore", dis)
	}
	if ann := annotations["joyrex2001.com/nightshift.schedule"]; ann != "" {
		return s.annotationToSchedule(ann)
	}
	if s.ForceSchedule != nil {
		return s.ForceSchedule, nil
	}
	return s.DefaultSchedule, nil
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
