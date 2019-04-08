package reply

import (
	"github.com/hirakiuc/gonta-app/log"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

// ClientConfig describe a config for slack.Client
type ClientConfig struct {
	AppToken string `envconfig:"APP_TOKEN" required:"true"`
}

// GetClient return a slack.Client instance.
func GetClient() (*slack.Client, error) {
	log := log.GetLogger()

	var conf ClientConfig
	if err := envconfig.Process("", &conf); err != nil {
		log.Error("Failed to load ClientConfig", zap.Error(err))
		return nil, err
	}

	return slack.New(conf.AppToken), nil
}
