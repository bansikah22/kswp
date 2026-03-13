package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	WebhookURL string
}

func (s *SlackNotifier) Send(message string) error {
	msg := slack.WebhookMessage{
		Text: message,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send slack notification, status code: %d", resp.StatusCode)
	}

	return nil
}
