package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event/data"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"go.uber.org/zap"
)

const (
	selectVersionBlockID     = "select-version"
	confirmDeploymentBlockID = "confirm-release"

	cancelVersion = "deny"
)

type Release struct {
	Base
}

func NewRelease(c *config.HandlerConfig, logger *zap.Logger) *Release {
	return &Release{
		Base: Base{
			config: c,
			logger: logger,
		},
	}
}

func (u *Release) needToRelease(msg string) bool {
	words := strings.Split(strings.TrimSpace(msg), " ")

	if len(words) == 0 {
		return false
	}

	return (strings.ToLower(words[1]) == "release")
}

/*
 * app mention.
 */
func (u *Release) ShowVersionChooser(e *slackevents.EventsAPIEvent) error {
	u.logger.Info("ShowVersionChooser start")

	ev, err := castAppMentionEvent(e)
	if err != nil {
		u.logger.Debug("Can't get AppMentionEvent...")

		return err
	}

	if !u.needToRelease(ev.Text) {
		u.logger.Debug("Release callback should not be invoked")
		// Ignore mention event

		return nil
	}

	u.logger.Info("Release flow start")

	textSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, "Please select *version*.", false, false),
		nil,
		nil,
	)

	selectMenu := slack.NewOptionsSelectBlockElement(
		slack.OptTypeExternal,
		slack.NewTextBlockObject(slack.PlainTextType, "Select version", false, false),
		"",
	)

	actionBlock := slack.NewActionBlock(selectVersionBlockID, selectMenu)

	fallbackText := slack.MsgOptionText("This client is not supported.", false)
	blocks := slack.MsgOptionBlocks(textSection, actionBlock)

	api := u.slackAPI()

	_, err = api.PostEphemeral(ev.Channel, ev.User, fallbackText, blocks)
	if err != nil {
		u.logger.Error("Failed to send an ephemeral message", zap.Error(err))

		return err
	}

	u.logger.Info("Sent a show versions message")

	return nil
}

func (u *Release) FetchVersions(e *data.ExternalDataRequest) ([]byte, error) {
	text := `{
	"options": [
		{
			"text": {
				"type": "plain_text",
				"text": "v1.0.0"
			},
			"value": "v1.0.0"
		},
		{
			"text": {
				"type": "plain_text",
				"text": "v1.1.0"
			},
			"value": "v1.1.0"
		},
		{
			"text": {
				"type": "plain_text",
				"text": "v1.1.1"
			},
			"value": "v1.1.1"
		}
	]
}`

	return []byte(text), nil
}

// actions.
func (u *Release) ConfirmRelease(e *slack.InteractionCallback) error {
	action := e.ActionCallback.BlockActions[0]
	version := action.SelectedOption.Value

	textSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("Could I deploy `%s`", version), false, false),
		nil,
		nil,
	)

	confirmButton := slack.NewButtonBlockElement(
		"",
		version,
		slack.NewTextBlockObject(slack.PlainTextType, "OK!", false, false),
	)
	confirmButton.WithStyle(slack.StylePrimary)

	denyButton := slack.NewButtonBlockElement(
		"",
		cancelVersion,
		slack.NewTextBlockObject(slack.PlainTextType, "Stop", false, false),
	)
	denyButton.WithStyle(slack.StyleDanger)

	actionBlock := slack.NewActionBlock(confirmDeploymentBlockID, confirmButton, denyButton)

	fallbackText := slack.MsgOptionText("This client is not supported.", false)
	blocks := slack.MsgOptionBlocks(textSection, actionBlock)

	replaceOriginal := slack.MsgOptionReplaceOriginal(e.ResponseURL)

	api := u.slackAPI()

	// nolint:dogsled
	_, _, _, err := api.SendMessage("", replaceOriginal, fallbackText, blocks)
	if err != nil {
		u.logger.Error("Failed to send a message", zap.Error(err))

		return err
	}

	return nil
}

/*
 * actions.
 */
func (u *Release) InvokeRelease(e *slack.InteractionCallback) error {
	action := e.ActionCallback.BlockActions[0]
	version := action.Value

	api := u.slackAPI()

	// Remove the original message to prevent double invoking this action.
	opt := slack.MsgOptionDeleteOriginal(e.ResponseURL)

	// nolint:dogsled
	_, _, _, err := api.SendMessage("", opt)
	if err != nil {
		u.logger.Error("Failed to delete the original message", zap.Error(err))

		return err
	}

	// Deploy should be cancelled if the version is the cancelVersion.
	if version == cancelVersion {
		msg := slack.MsgOptionText(
			fmt.Sprintf("<@%s> Cancelled!", e.User.ID),
			false,
		)

		_, _, err := api.PostMessage(e.Channel.ID, msg)
		if err != nil {
			u.logger.Error("Failed to send a cancel message")

			return err
		}

		return nil
	}

	startMsg := slack.MsgOptionText(
		fmt.Sprintf("<@%s> OK, I'll deploy `%s`.", e.User.ID, version),
		false,
	)

	_, _, err = api.PostMessage(e.Channel.ID, startMsg)
	if err != nil {
		u.logger.Error("Failed to send a start message", zap.Error(err))

		return err
	}

	ch := make(chan error)

	// Dispatch deploy process
	go u.deploy(ch, e, version)

	u.logger.Info("Waiting for the deployment", zap.String("version", version))

	return <-ch
}

func (u *Release) deploy(ch chan error, e *slack.InteractionCallback, version string) {
	// Wait the deployment
	// nolint:gomnd
	time.Sleep(3 * time.Second)

	api := u.slackAPI()

	u.logger.Info("Start deploying the version", zap.String("version", version))

	// u.deploy(version)

	endMsg := slack.MsgOptionText(
		fmt.Sprintf("`%s` deployed", version),
		false,
	)

	_, _, err := api.PostMessage(e.Channel.ID, endMsg)
	if err != nil {
		u.logger.Error("Failed to send a complete message", zap.Error(err))

		ch <- err

		return
	}

	u.logger.Info("Deployed the version", zap.String("version", version))

	ch <- nil
}
