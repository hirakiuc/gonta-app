package config

import (
	"github.com/kelseyhightower/envconfig"
)

type HandlerConfig struct {
	BotAccessToken string
}

type SlackConfig struct {
	BotAccessToken    string `envconfig:"BOT_USER_ACCESS_TOKEN" required:"true"`
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`
	SigningSecret     string `envconfig:"SLACK_SIGNING_SECRET" required:"true"`
}

type Config struct {
	Slack *SlackConfig
}

func NewConfig() *Config {
	return &Config{
		Slack: &SlackConfig{},
	}
}

func (c *Config) Load() error {
	return envconfig.Process("", c)
}

func (c *Config) HandlerConfig() *HandlerConfig {
	return &HandlerConfig{
		BotAccessToken: c.Slack.BotAccessToken,
	}
}
