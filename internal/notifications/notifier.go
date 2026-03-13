package notifications

import (
	"fmt"

	"github.com/bansikah22/kswp/internal/config"
)

type Notifier interface {
	Send(message string) error
}

var NewNotifier = func(config config.NotificationConfig) (Notifier, error) {
	switch config.Type {
	case "slack":
		return &SlackNotifier{WebhookURL: config.WebhookURL}, nil
	case "email":
		return &EmailNotifier{
			From:     config.Email.From,
			To:       config.Email.To,
			Password: config.Email.Password,
			Host:     config.Email.Host,
			Port:     config.Email.Port,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported notifier type: %s", config.Type)
	}
}
