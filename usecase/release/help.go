package release

import (
	"strings"

	"github.com/hirakiuc/gonta-app/usecase"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Help struct {
	usecase.Base
}

func NewHelp(r *Release) *Help {
	return &Help{
		Base: r.Base,
	}
}

func (h *Help) Show(e *slackevents.AppMentionEvent) error {
	base := `
1. @BOT_ID release
    -> Show usage
2. @BOT_ID release help
    -> Show usage
3. @BOT_ID release [repo]
    -> Choose version & deploy
4. @BOT_ID release [repo] [version]
    -> Deploy the version
`

	// nolint:godox
	// TODO Replace BOT_ID
	msg := strings.ReplaceAll(base, "@BOT_ID", "`<@bot_id>`")

	section := slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, msg, false, false),
		nil,
		nil,
	)

	block := slack.MsgOptionBlocks(section)

	api := h.SlackAPI()

	// nolint:dogsled
	_, _, _, err := api.SendMessage(e.Channel, block)
	if err != nil {
		h.Logger.Error("Failed to send a message", zap.Error(err))

		return err
	}

	return nil
}
