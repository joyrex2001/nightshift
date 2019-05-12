package trigger

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
	// TODO: call Webhook/handle timeout
	return nil
}
