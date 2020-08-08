package server

import (
	"encoding/json"

	"github.com/hirakiuc/gonta-app/event/data"
)

func ParseExternalDataRequest(v []byte) (*data.ExternalDataRequest, error) {
	var req data.ExternalDataRequest

	err := json.Unmarshal(v, &req)

	return &req, err
}
