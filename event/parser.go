package event

import (
	"bytes"
	"encoding/json"
	"strings"

	scan "github.com/mattn/go-scan"
)

// Parser describe a parser instance to parse the slack events.
type Parser struct {
}

// NewParser return a pointer to new parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// GetType extract type value in the json string.
func (p *Parser) GetType(jsonStr string) (*string, error) {
	reader := strings.NewReader(jsonStr)

	var typeStr string
	err := scan.ScanJSON(reader, "/type", &typeStr)
	if err != nil {
		return nil, err
	}
	return &typeStr, nil
}

// GetToken extract token value in the json string.
func (p *Parser) GetToken(jsonStr string) (*string, error) {
	reader := strings.NewReader(jsonStr)

	var token string
	err := scan.ScanJSON(reader, "/token", &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// ParseCallbackEvent parse the event as a CallbackEvent
func (p *Parser) ParseCallbackEvent(jsonStr string) (*CallbackEvent, error) {
	e := CallbackEvent{}
	if err := json.Unmarshal([]byte(jsonStr), &e); err != nil {
		return nil, err
	}

	return &e, nil
}

// ParseURLVerificationEvent parse the event as a URLVerificationEvent
func (p *Parser) ParseURLVerificationEvent(jsonStr string) (*URLVerificationEvent, error) {
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
