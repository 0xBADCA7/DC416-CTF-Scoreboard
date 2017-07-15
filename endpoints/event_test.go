package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventInfo(test *testing.T) {
	handler := NewEventInfoHandler("test event")
	server := httptest.NewServer(handler)
	defer server.Close()

	// This endpoint will always serve JSON containing the event name
	response, err := http.Get(server.URL())
	if err != nil {
		test.Error(err)
	}
	if response.StatusCode != 200 {
		test.Errorf("Expected the request to succeed. Got status code %d\n", response.StatusCode)
	}
	data := map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()
	err = decoder.Decode(&data)
	if err != nil {
		test.Error(err)
	}
	name, found := data["name"]
	if !found {
		test.Errorf("Expected response to have a 'name' field, but it does not\n")
	}
	if name != "test event" {
		test.Errorf("Expected the event name to be 'test event' but got '%s'\n", name) 
	}
}
