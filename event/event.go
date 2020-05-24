package event

import "encoding/json"

/*
token:
challenge:
type: "url_verification"

---

token:
type:

team_id:
api_app_id:
event: ...
event_id:
event_time:
authed_users: [
	...
	]
*/

// URLVerificationEvent describe a url verification event.
type URLVerificationEvent struct {
	Challenge string `json:"challenge" required:"true"`
	Token     string `json:"token" required:"true"`
	Type      string `json:"type" required:"true"`
}

// CallbackEvent describe a event which sent from slack.
type CallbackEvent struct {
	APIAppID    string          `json:"api_app_id" required:"true"`
	AuthedUsers []string        `json:"authed_users" required:"true"`
	Event       json.RawMessage `json:"event" required:"true"`
	EventID     string          `json:"event_id" required:"true"`
	EventTime   int64           `json:"event_time" required:"true"`
	TeamID      string          `json:"team_id" required:"true"`
	Token       string          `json:"token" required:"true"`
	Type        string          `json:"type" required:"true"`
}

/** A kind of event payload **/

// AppMentionEvent describe the actual event from slack.
type AppMentionEvent struct {
	Channel     string `json:"channel" required:"true"`
	ClientMsgID string `json:"client_msg_id" required:"true"`
	EventTS     string `json:"event_ts" required:"true"`
	Text        string `json:"text" required:"true"`
	TS          string `json:"ts" required:"true"`
	Type        string `json:"type" required:"true"`
	User        string `json:"user" required:"true"`
}
