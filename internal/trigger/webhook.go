package trigger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
)

// WebhookTrigger is the object that implements http based triggers.
type WebhookTrigger struct {
	config Config
}

func init() {
	RegisterModule("webhook", NewWebhookTrigger)
}

// NewWebhookTrigger will instantiate a new WebhookTrigger object.
func NewWebhookTrigger() (Trigger, error) {
	return &WebhookTrigger{config: Config{}}, nil
}

// SetConfig will set the generic configuration for this trigger.
func (s *WebhookTrigger) SetConfig(cfg Config) {
	s.config = cfg
}

// GetConfig will return the config applied for this trigger.
func (s *WebhookTrigger) GetConfig() Config {
	return s.config
}

// Execute will trigger the webhook.
func (s *WebhookTrigger) Execute() error {
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("error webhook; status=%s(%d)", resp.Status, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	glog.V(5).Infof("url: %s, status: %s, body: %s", s.config.Settings["url"], resp.Status, body)
	return nil
}

// newClient will create a new http.Client object with the correct settings, as
// reflected in the config.
func (s *WebhookTrigger) newClient() (*http.Client, error) {
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
func (s *WebhookTrigger) newRequest() (*http.Request, error) {
	method := s.getMethod()
	url, err := s.getUrl()
	if err != nil {
		return nil, err
	}
	body, err := s.getBody()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	headers, err := s.getHeaders()
	if err != nil {
		return nil, err
	}
	for headr, val := range headers {
		req.Header.Set(headr, val)
	}
	return req, nil
}

func (s *WebhookTrigger) getUrl() (string, error) {
	url := strings.TrimSpace(s.config.Settings["url"])
	if url == "" {
		return "", fmt.Errorf("no url specified")
	}
	return url, nil
}

// getBody will process the configured body, and return an io.ReadWriter for
// that body.
func (s *WebhookTrigger) getBody() (io.ReadWriter, error) {
	buf := new(bytes.Buffer)
	body, err := RenderTemplate(s.config.Settings["body"], s.config.Settings)
	if err != nil {
		return buf, err
	}
	if body != "" {
		buf.WriteString(body)
	}
	return buf, nil
}

// getTimeout will return a time.Duration for the configured timeout. If no
// timeout has been configured it will use a default timeout instead.
func (s *WebhookTrigger) getTimeout() (time.Duration, error) {
	to := s.config.Settings["timeout"]
	if to == "" {
		to = "300ms"
	}
	return time.ParseDuration(to)
}

// getMethod will return the method as configured. If no method is set, it will
// default to GET if no body is configured, or POST if a body has been
// configured.
func (s *WebhookTrigger) getMethod() string {
	method := strings.ToUpper(s.config.Settings["method"])
	if method == "" {
		method = "GET"
		if s.config.Settings["body"] != "" {
			method = "POST"
		}
	}
	return method
}

// getHeaders will parse the headers configuration, and return a map containing
// the headers and its' values.
func (s *WebhookTrigger) getHeaders() (map[string]string, error) {
	headers := map[string]string{}
	chdrs := strings.Split(strings.Replace(s.config.Settings["headers"], "\r\n", "\n", -1), "\n")
	for _, header := range chdrs {
		if header == "" {
			continue
		}
		flds := strings.Split(header, ":")
		if len(flds) != 2 {
			return nil, fmt.Errorf("invalid header specified '%s'", header)
		}
		flds[0] = strings.TrimSpace(flds[0])
		flds[1] = strings.TrimSpace(flds[1])
		if flds[0] == "" || flds[1] == "" {
			continue
		}
		val, err := RenderTemplate(flds[1], s.config.Settings)
		if err != nil {
			return headers, err
		}
		headers[flds[0]] = val
	}
	return headers, nil
}
