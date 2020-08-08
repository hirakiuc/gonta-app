package data

import (
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
	Options        []*Element `json:"option"`
	InitialOption  *Element   `json:"initial_option,omitempty"`
	MinQueryLength int64      `json:"min_query_length"` // default:3
}

type OptionGroups struct {
	Groups []*Options `json:"option_groups"`
}

/**
 *
 * NOTE: https://api.slack.com/reference/block-kit/block-elements#external_select
 */
type Element struct {
	Text        string                 `json:"type"`
	Placeholder *slack.TextBlockObject `json:"placeholder"`
	ActionID    string                 `json:"action_id"`
}
