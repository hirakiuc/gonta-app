package release

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hirakiuc/gonta-app/config"
	"github.com/hirakiuc/gonta-app/event/data"
	"github.com/hirakiuc/gonta-app/usecase"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

const (
	SelectVersionBlockID     = "select-version"
	ConfirmDeploymentBlockID = "confirm-release"

	VersionChooserActionID = "select-version"
	CancelVersion          = "deny"

	SubCommandHelp = "help"

	MinQueryLength = 2
)

type Release struct {
	usecase.Base
}

func New(c *config.HandlerConfig, logger *zap.Logger) *Release {
	return &Release{
		Base: usecase.Base{
			Config: c,
			Logger: logger,
		},
	}
}

func (u *Release) Start(e *slackevents.EventsAPIEvent) error {
	u.Logger.Info("Invoke release flow.")

	ev, err := u.CastAppMentionEvent(e)
	if err != nil {
		u.Logger.Error("Can't get AppMentionEvent...", zap.Error(err))

		return err
	}

	cmd := u.ParseAsCommand(ev.Text, "release")
	if cmd == nil {
		u.Logger.Debug("Ignroe this event")

		return nil
	}

	// Usage:
	// 1. @gonta release
	//   -> show usage
	// 2. @gonta release help
	//   -> show usage
	// 3. @gonta release [repo]
	//   -> confirm the repo & fetch tags
	// 4. @gonta release [repo] [version]
	//   -> confirm the repo & version & deploy

	switch len(cmd.Args) {
	case 0:
		// @gonta release
		h := NewHelp(u)

		return h.Show(ev)

	case 1:
		if cmd.Args[0] == SubCommandHelp {
			// @gonta release help
			h := NewHelp(u)

			return h.Show(ev)
		}

		// @gonta release [repo]
		v := NewVersionChooser(u)

		return v.Show(ev, cmd.Args[0])

	case 2: // nolint:gomnd
		// @gonta release [repo] [version]
		c := NewConfirm(u)

		return c.Confirm(ev.Channel, cmd.Args[0], cmd.Args[1])

	default:
		// @gonta release A B C....
		h := NewHelp(u)

		return h.Show(ev)
	}
}

// This actionID will be used in version chooser.
func actionIDWithRepo(repo string) string {
	return fmt.Sprintf("%s:%s", VersionChooserActionID, repo)
}

// This method is used in External Data Fetcher to get the repository name from the ActionID.
func parseActionID(actionID string) string {
	if !strings.HasPrefix(actionID, VersionChooserActionID+":") {
		return ""
	}

	pos := len(VersionChooserActionID + ":")

	return actionID[pos:]
}

// Parse actionID and extract repo.
func parseSelectedVersion(actionID string) (string, string) {
	parts := strings.Split(actionID, ":")

	repo := parts[0]
	version := strings.Join(parts[1:], ":")

	return repo, version
}

func (u *Release) FetchVersions(e *data.ExternalDataRequest) ([]byte, error) {
	repo := parseActionID(e.ActionID)
	if len(repo) == 0 {
		u.Logger.Debug("No repo found from the actionID", zap.String("actionID", e.ActionID))

		return []byte(`{"options":[]}`), nil
	}

	u.Logger.Debug(
		"Repository",
		zap.String("ActionID", e.ActionID),
		zap.String("repo", repo),
	)

	// Fetch versions
	f := NewVersionFetcher(u)

	versions, err := f.Fetch(repo, e.Value)
	if err != nil {
		u.Logger.Error("Failed to fetch versions", zap.Error(err))

		return []byte(`{"options":[]}`), err
	}

	options := data.NewOptions()
	options.AddVersionsWithRepo(repo, versions)

	bytes, err := json.Marshal(options)
	if err != nil {
		u.Logger.Error("Failed to generate json", zap.Error(err))

		return []byte(`{"options":[]}`), err
	}

	u.Logger.Debug("response", zap.String("json", string(bytes)))

	return bytes, nil
}

// actions.
func (u *Release) ConfirmRelease(e *slack.InteractionCallback) error {
	action := e.ActionCallback.BlockActions[0]
	value := action.SelectedOption.Value

	repo, version := parseSelectedVersion(value)

	c := NewConfirm(u)

	return c.ConfirmFromCallback(e, repo, version)
}

/*
 * actions.
 */
func (u *Release) InvokeRelease(e *slack.InteractionCallback) error {
	action := e.ActionCallback.BlockActions[0]

	repo := parseActionID(action.ActionID)

	version := action.Value

	d := NewDeploy(u, repo, version)

	return d.Start(e)
}
