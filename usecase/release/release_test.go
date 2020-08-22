package release

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseActionID(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		Text            string
		ExpectedRepo    string
		ExpectedVersion string
	}{
		{
			Text:            "repo:v1.0.0",
			ExpectedRepo:    "repo",
			ExpectedVersion: "v1.0.0",
		},
	}

	for _, testcase := range cases {
		repo, version := parseSelectedVersion(testcase.Text)

		assert.Equal(testcase.ExpectedRepo, repo)
		assert.Equal(testcase.ExpectedVersion, version)
	}
}
