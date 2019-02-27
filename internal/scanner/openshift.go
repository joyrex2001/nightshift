package scanner

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

func (s *OpenShiftScanner) GetSchedule() ([]schedule.Schedule, error) {
	if s.kubernetes == nil {
		return nil, fmt.Errorf("unable to connect to kubernetes")
	}

	return nil, nil
}
