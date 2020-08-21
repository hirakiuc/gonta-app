package usecase

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	Separator = " "
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

type Base struct {
	Config *config.HandlerConfig
	Logger *zap.Logger
}

type Command struct {
	Name string
	Args []string
}

func (b *Base) SlackAPI() *slack.Client {
	// nolint:gosec
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return slack.New(b.Config.BotAccessToken, slack.OptionHTTPClient(c))
}

func (b *Base) CastAppMentionEvent(e *slackevents.EventsAPIEvent) (*slackevents.AppMentionEvent, error) {
	innerEvent := e.InnerEvent

	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		return ev, nil
	default:
		return nil, fmt.Errorf("%w", ErrUnexpectedEventType)
	}
}

func findWordInArray(ary []string, word string) (int, bool) {
	target := strings.ToLower(word)

	for idx, v := range ary {
		if strings.ToLower(v) == target {
			return idx, true
		}
	}

	return 0, false
}

/*
 * Extract the origin word from auto-linked text by slack.
 *
 * NOTE: extract the text from link markdown, "select-version:<http://github.com/gonta-app|github.com/gonta-app>"
 */
func (b *Base) parseAutoLinkWord(text string) string {
	if !strings.HasPrefix(text, "<http://") || !strings.HasSuffix(text, ">") {
		return text
	}

	if 0 > strings.Index(text, "|") {
		return text
	}

	pos := strings.Index(text, "|")

	start := pos + 1
	end := len(text) - 1

	if start >= end {
		return text
	}

	return text[start:end]
}

func (b *Base) ParseAsCommand(text string, startFrom string) *Command {
	parts := strings.Split(text, Separator)

	values := []string{}

	for _, v := range parts {
		word := strings.TrimSpace(v)
		if len(word) == 0 {
			continue
		}

		values = append(values, word)
	}

	pos, ok := findWordInArray(values, startFrom)
	if !ok {
		return nil
	}

	words := values[pos:]

	if len(words) == 0 {
		return nil
	}

	args := make([]string, len(words[1:]))
	for i, v := range words[1:] {
		args[i] = b.parseAutoLinkWord(v)
	}

	return &Command{
		Name: words[0],
		Args: args,
	}
}
