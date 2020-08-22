package release

import (
	"fmt"

	"github.com/hirakiuc/gonta-app/usecase"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Confirm struct {
	usecase.Base
}

func NewConfirm(u *Release) *Confirm {
	return &Confirm{
		Base: u.Base,
	}
}

func (c *Confirm) ConfirmFromCallback(e *slack.InteractionCallback, repo string, version string) error {
	return c.showConfirm(e.Channel.ID, repo, version, e.ResponseURL)
}

func (c *Confirm) Confirm(channelID string, repo string, version string) error {
	return c.showConfirm(channelID, repo, version, "")
}

func (c *Confirm) showConfirm(channelID string, repo string, version string, responseURL string) error {
	msg := fmt.Sprintf("We're goging to deploy version `%s` of `%s`.\nAre you sure?", version, repo)
	textSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, msg, false, false),
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
		CancelVersion,
		slack.NewTextBlockObject(slack.PlainTextType, "Stop", false, false),
	)
	denyButton.WithStyle(slack.StyleDanger)

	actionBlock := slack.NewActionBlock(ConfirmDeploymentBlockID, confirmButton, denyButton)

	// Build options for SendMessage
	opts := []slack.MsgOption{}

	fallbackText := slack.MsgOptionText("This client is not supported.", false)
	opts = append(opts, fallbackText)

	blocks := slack.MsgOptionBlocks(textSection, actionBlock)
	opts = append(opts, blocks)

	if len(responseURL) > 0 {
		replaceOriginal := slack.MsgOptionReplaceOriginal(responseURL)
		opts = append(opts, replaceOriginal)
	}

	api := c.SlackAPI()

	// nolint:dogsled
	_, _, _, err := api.SendMessage(channelID, opts...)
	if err != nil {
		c.Logger.Error("Failed to send a message", zap.Error(err))

		return err
	}

	return nil
}
