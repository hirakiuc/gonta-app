package event

import (
	"bytes"
	"encoding/json"

	scan "github.com/mattn/go-scan"
)

// Parser describe a parser instance to parse the slack events.
type EventParser struct {
}

type EventParseResult struct {
	JSON  string
	Type  string
	Token string
}

// NewEventParser return a pointer to new parser instance.
func NewEventParser() *EventParser {
	return &EventParser{}
}

// ParseCallbackEvent parse the event as a CallbackEvent.
func (p *EventParser) ParseCallbackEvent(jsonStr string) (*CallbackEvent, error) {
	e := CallbackEvent{}
	if err := json.Unmarshal([]byte(jsonStr), &e); err != nil {
		return nil, err
	}

	return &e, nil
}

// ParseURLVerificationEvent parse the event as a URLVerificationEvent.
func (p *EventParser) ParseURLVerificationEvent(jsonStr string) (*URLVerificationEvent, error) {
	e := URLVerificationEvent{}
	if err := json.Unmarshal([]byte(jsonStr), &e); err != nil {
		return nil, err
	}

	return &e, nil
}

// GetEventType extract the type value in the event json.
func (e *CallbackEvent) GetEventType() (*string, error) {
	reader := bytes.NewReader(e.Event)

	var token string
	err := scan.ScanJSON(reader, "/type", &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// ParseAppMentionEvent parse the event as a AppMentionEvent.
func (e *CallbackEvent) ParseAppMentionEvent() (*AppMentionEvent, error) {
	v := AppMentionEvent{}
	if err := json.Unmarshal(e.Event, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
