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
			Text: "release opt1 opt2",
			Cmd:  nil,
		},
		{
			Text: "@gonta release opt1   opt2  ",
			Cmd: &Command{
				Name: "release",
				Args: []string{"opt1", "opt2"},
			},
		},
	}

	base := &Base{}

	for _, testcase := range cases {
		result := base.ParseAsCommand(testcase.Text, "@gonta")
		assert.Equal(testcase.Cmd, result)
	}
}
