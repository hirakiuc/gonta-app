package release

import (
	"fmt"
	"time"

	"github.com/hirakiuc/gonta-app/usecase"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Deploy struct {
	Repo    string
	Version string

	usecase.Base
}

func NewDeploy(r *Release, repo string, version string) *Deploy {
	return &Deploy{
		Repo:    repo,
		Version: version,
		Base:    r.Base,
	}
}

func (d *Deploy) Start(e *slack.InteractionCallback) error {
	api := d.SlackAPI()

	// Remove the original message to prevent double invoking this action.
	opt := slack.MsgOptionDeleteOriginal(e.ResponseURL)

	// nolint:dogsled
	_, _, _, err := api.SendMessage("", opt)
	if err != nil {
		d.Logger.Error("Failed to delete the original message", zap.Error(err))

		return err
	}

	// Deploy should be cancelled if the version is the cancelVersion.
	if d.isCancelRequest() {
		return d.cancel(e, api)
	}

	startMsg := slack.MsgOptionText(
		fmt.Sprintf("<@%s> OK, I'll deploy `%s`.", e.User.ID, d.Version),
		false,
	)

	_, _, err = api.PostMessage(e.Channel.ID, startMsg)
	if err != nil {
		d.Logger.Error("Failed to send a start message", zap.Error(err))

		return err
	}

	ch := make(chan error)

	// Dispatch deploy process
	go d.deploy(ch, e, d.Version)

	d.Logger.Info("Waiting for the deployment", zap.String("version", d.Version))

	return <-ch
}

func (d *Deploy) isCancelRequest() bool {
	return d.Version == CancelVersion
}

func (d *Deploy) cancel(e *slack.InteractionCallback, api *slack.Client) error {
	msg := slack.MsgOptionText(
		fmt.Sprintf("<@%s> Cancelled!", e.User.ID),
		false,
	)

	_, _, err := api.PostMessage(e.Channel.ID, msg)
	if err != nil {
		d.Logger.Error("Failed to send a cancel message")

		return err
	}

	return nil
}

func (d *Deploy) deploy(ch chan error, e *slack.InteractionCallback, version string) {
	// Wait the deployment
	// nolint:gomnd
	time.Sleep(3 * time.Second)

	api := d.SlackAPI()

	d.Logger.Info("Start deploying the version", zap.String("version", version))

	// u.deploy(version)

	endMsg := slack.MsgOptionText(
		fmt.Sprintf("`%s` deployed", version),
		false,
	)

	_, _, err := api.PostMessage(e.Channel.ID, endMsg)
	if err != nil {
		d.Logger.Error("Failed to send a complete message", zap.Error(err))

		ch <- err

		return
	}

	d.Logger.Info("Deployed the version", zap.String("version", version))

	ch <- nil
}
