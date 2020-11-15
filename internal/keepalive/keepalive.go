package keepalive

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/joyrex2001/nightshift/internal/scanner"
)

// KeepAlive defines the public interface.
type KeepAlive interface {
	SetConfig(Config)
	GetConfig() Config
	Execute([]*scanner.Object) error
}

// Config is the configuration for this keepalive, and contains a hashmap with
// generic settings. The key for each value should be lowercased always.
type Config struct {
	Id       string            `json:"id"`
	Settings map[string]string `json:"settings"`
}

// WebhookKeepAlive is the object that implements http based triggers.
type WebhookKeepAlive struct {
	config Config
}

// New will return a KeepAlive object.
func New() (KeepAlive, error) {
	return &WebhookKeepAlive{config: Config{}}, nil
}

// SetConfig will set the generic configuration for this trigger.
func (s *WebhookKeepAlive) SetConfig(cfg Config) {
	s.config = cfg
}

// GetConfig will return the config applied for this trigger.
func (s *WebhookKeepAlive) GetConfig() Config {
	return s.config
}

// Execute will trigger the webhook.
func (s *WebhookKeepAlive) Execute(objs []*scanner.Object) error {
	cli, err := s.newClient()
	if err != nil {
		return err
	}
	req, err := s.newRequest()
	if err != nil {
		return err
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("error webhook; status=%s(%d)", resp.Status, resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	glog.V(5).Infof("url: %s, status: %s, body: %s", s.config.Settings["url"], resp.Status, body)
	return nil
}

// newClient will create a new http.Client object with the correct settings, as
// reflected in the config.
func (s *WebhookKeepAlive) newClient() (*http.Client, error) {
	timeout, err := s.getTimeout()
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Timeout: timeout,
	}, nil
}

// newRequest will create a http.Request for the configured url, body and
// method.
func (s *WebhookKeepAlive) newRequest() (*http.Request, error) {
	url := strings.TrimSpace(s.config.Settings["url"])
	buf := new(bytes.Buffer)
	req, err := http.NewRequest("GET", url, buf)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// getTimeout will return a time.Duration for the configured timeout. If no
// timeout has been configured it will use a default timeout instead.
func (s *WebhookKeepAlive) getTimeout() (time.Duration, error) {
	to := s.config.Settings["timeout"]
	if to == "" {
		to = "300ms"
	}
	return time.ParseDuration(to)
}
