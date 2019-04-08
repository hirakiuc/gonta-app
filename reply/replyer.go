package reply

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/event"
)

type Replyer interface {
	Reply(w http.ResponseWriter, msg *event.CallbackEvent)
}
