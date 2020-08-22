package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAsCommand(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		Text string
		Cmd  *Command
	}{
		{
			Text: "@gonta ",
			Cmd:  nil,
		},
		{
			Text: " @gonta release ",
			Cmd: &Command{
				Name: "release",
				Args: []string{},
			},
		},
		{
			Text: "@gonta release opt1   opt2  ",
			Cmd: &Command{
				Name: "release",
				Args: []string{"opt1", "opt2"},
			},
		},
		{
			Text: "@gonta release <http://github.com/hirakiuc/gonta-app|github.com/hirakiuc/gonta-app>",
			Cmd: &Command{
				Name: "release",
				Args: []string{"github.com/hirakiuc/gonta-app"},
			},
		},
	}

	base := &Base{}

	for _, testcase := range cases {
		result := base.ParseAsCommand(testcase.Text, "release")
		assert.Equal(testcase.Cmd, result)
	}
}

func TestParseAutoLinkWord(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		Text     string
		Expected string
	}{
		{
			Text:     "word",
			Expected: "word",
		},
		{
			Text:     "<http://github.com/hirakiuc/gonta-app|github.com/hirakiuc/gonta-app>",
			Expected: "github.com/hirakiuc/gonta-app",
		},
		{
			Text:     "<http://a|a>",
			Expected: "a",
		},
	}

	base := &Base{}

	for _, testcase := range cases {
		result := base.parseAutoLinkWord(testcase.Text)

		assert.Equal(testcase.Expected, result)
	}
}
