package event

import (
	"strings"

	"github.com/hirakiuc/gonta-app/handler"
	scan "github.com/mattn/go-scan"
)

// BodyParser describe a parser instance to parse request body
type BodyParser struct {
}

type BodyParseResult struct {
	JSON  string
	Type  string
	Token string
}

// NewBodyParser return a pointer to new parser instance.
func NewBodyParser() *BodyParser {
	return &BodyParser{}
}

func (p *BodyParser) Parse(json string) (*BodyParseResult, error) {
	result := BodyParseResult{JSON: json}

	var err error

	result.Type, err = p.getType(json)
	if err != nil {
		return nil, err
	}

	result.Token, err = p.getToken(json)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetType extract type value in the json string.
func (p *BodyParser) getType(jsonStr string) (string, error) {
	reader := strings.NewReader(jsonStr)

	var typeStr string

	err := scan.ScanJSON(reader, "/type", &typeStr)
	if err != nil {
		return "", err
	}

	return typeStr, nil
}

// GetToken extract token value in the json string.
func (p *BodyParser) getToken(jsonStr string) (string, error) {
	reader := strings.NewReader(jsonStr)

	var token string

	err := scan.ScanJSON(reader, "/token", &token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (result *BodyParseResult) EventHandler() (handler.Handler, error) {

}
