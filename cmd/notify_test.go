package cmd

import (
	"bytes"
	"testing"

	"github.com/bansikah22/kswp/internal/notifications"
	"github.com/stretchr/testify/assert"

	"github.com/bansikah22/kswp/internal/config"
)

func TestNotifySlackCmd(t *testing.T) {
	// Redirect stdout to a buffer
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	// Create a mock notifier
	mockNotifier := &notifications.MockNotifier{}
	oldNewNotifier := notifications.NewNotifier
	defer func() { notifications.NewNotifier = oldNewNotifier }()
	notifications.NewNotifier = func(config config.NotificationConfig) (notifications.Notifier, error) {
		return mockNotifier, nil
	}

	// Execute the command
	rootCmd.SetArgs([]string{"notify", "slack", "--webhook-url", "test", "--message", "hello"})
	rootCmd.Execute()

	// Check the output
	assert.Contains(t, buf.String(), "Email notification sent successfully.")
}

func TestNotifyEmailCmd(t *testing.T) {
	// Redirect stdout to a buffer
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	// Create a mock notifier
	mockNotifier := &notifications.MockNotifier{}
	oldNewNotifier := notifications.NewNotifier
	defer func() { notifications.NewNotifier = oldNewNotifier }()
	notifications.NewNotifier = func(config config.NotificationConfig) (notifications.Notifier, error) {
		return mockNotifier, nil
	}

	// Execute the command
	rootCmd.SetArgs([]string{"notify", "email", "--message", "hello"})
	rootCmd.Execute()

	// Check the output
	assert.Contains(t, buf.String(), "Email notification sent successfully.")
}
