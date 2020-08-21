package data

import (
	"fmt"

	"github.com/slack-go/slack"
)

type ExternalDataRequest struct {
	Type      string           `json:"type"`
	User      *slack.User      `json:"user"`
	Container *slack.Container `json:"container"`
	APIAppID  string           `json:"api_app_id"`
	Token     string           `json:"token"`
	ActionID  string           `json:"action_id"`
	BlockID   string           `json:"block_id"`
	Value     string           `json:"value"`
	Team      *slack.Team      `json:"team"`
	Channel   *slack.Channel   `json:"channel"`
}

type Options struct {
	Options       []*Option `json:"options"`
	InitialOption *Option   `json:"initial_option,omitempty"`
}

func NewOptions() *Options {
	return &Options{
		Options:       []*Option{},
		InitialOption: nil,
	}
}

func (o *Options) AddVersionsWithRepo(repo string, versions []string) {
	for _, version := range versions {
		v := fmt.Sprintf("%s:%s", repo, version)

		o.AddOption(v, v)
	}
}

func (o *Options) AddOption(text string, value string) {
	opt := NewOption(text, value)

	o.Options = append(o.Options, opt)
}

type OptionGroups struct {
	Groups []*Options `json:"option_groups"`
}

/**
 *
 * NOTE: https://api.slack.com/reference/block-kit/block-elements#external_select
 */
type Option struct {
	Text  *TextObject `json:"text"`
	Value string      `json:"value"`
}

type TextObject struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewOption(text string, v string) *Option {
	return &Option{
		Text: &TextObject{
			Type: slack.PlainTextType,
			Text: text,
		},
		Value: v,
	}
}
