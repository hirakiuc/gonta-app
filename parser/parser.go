package parser

import (
	"encoding/json"

	"github.com/hirakiuc/gonta-app/model"
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
func (p *EventParser) ParseCallbackEvent(jsonStr string) (*model.CallbackEvent, error) {
	e := model.CallbackEvent{}
	if err := json.Unmarshal([]byte(jsonStr), &e); err != nil {
		return nil, err
	}

	return &e, nil
}

// ParseURLVerificationEvent parse the event as a URLVerificationEvent.
func (p *EventParser) ParseURLVerificationEvent(jsonStr string) (*model.URLVerificationEvent, error) {
	e := model.URLVerificationEvent{}
	if err := json.Unmarshal([]byte(jsonStr), &e); err != nil {
		return nil, err
	}

	return &e, nil
}
