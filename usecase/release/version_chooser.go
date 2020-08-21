package release

import (
	"github.com/hirakiuc/gonta-app/usecase"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type VersionChooser struct {
	usecase.Base
}

func NewVersionChooser(release *Release) *VersionChooser {
	return &VersionChooser{
		Base: release.Base,
	}
}

func (v *VersionChooser) Show(e *slackevents.AppMentionEvent, repo string) error {
	v.Logger.Info("ShowVersionChooser start")

	textSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, "Please select *version*.", false, false),
		nil,
		nil,
	)

	actionID := actionIDWithRepo(repo)
	v.Logger.Info(
		"ActionID in chooser",
		zap.String("ActionID", actionID),
		zap.String("repo", repo),
	)

	selectMenu := slack.NewOptionsSelectBlockElement(
		slack.OptTypeExternal,
		slack.NewTextBlockObject(slack.PlainTextType, "Select version", false, false),
		actionID,
	)

	actionBlock := slack.NewActionBlock(SelectVersionBlockID, selectMenu)

	fallbackText := slack.MsgOptionText("This client is not supported.", false)
	blocks := slack.MsgOptionBlocks(textSection, actionBlock)

	api := v.SlackAPI()

	_, err := api.PostEphemeral(e.Channel, e.User, fallbackText, blocks)
	if err != nil {
		v.Logger.Error("Failed to send an ephemeral message", zap.Error(err))

		return err
	}

	v.Logger.Info("Sent a show versions message")

	return nil
}
