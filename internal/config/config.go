package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Notifications []NotificationConfig `yaml:"notifications"`
}

type NotificationConfig struct {
	Type       string       `yaml:"type"`
	WebhookURL string       `yaml:"webhook_url,omitempty"`
	Email      *EmailConfig `yaml:"email,omitempty"`
}

type EmailConfig struct {
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
	Password string   `yaml:"password"`
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(home, ".kswp", "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
