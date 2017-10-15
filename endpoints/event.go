package endpoints

import (
	"encoding/json"
	"net/http"
)

// EventInfoHandler handles requests to the /event endpoint, serving descriptive information about the CTF event.
type EventInfoHandler struct {
	ctfName string
}

// EventInfo contains public data that can be encoded into the response to a request for event info.
type EventInfo struct {
	Name string `json:"name"`
}

// NewEventInfoHandler constructs a new EventInfoHandler.
func NewEventInfoHandler(ctfName string) EventInfoHandler {
	return EventInfoHandler{
		ctfName,
	}
}

func (self EventInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := EventInfo{self.ctfName}
	encoder := json.NewEncoder(w)
	encoder.Encode(&response)
}
