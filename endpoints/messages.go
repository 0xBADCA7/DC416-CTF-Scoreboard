package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// MessagesHandler handles requests to retrieve messages left for CTF participants by admins.
type MessagesHandler struct {
	messages models.MessageModel
}

type messagesResponse struct {
	Messages []models.Message `json:"messages"`
}

// NewMessagesHandler constructs a new MessagesHandler with the capability to access messages.
func NewMessagesHandler(messages models.MessageModel) MessagesHandler {
	return MessagesHandler{
		messages,
	}
}

func (self MessagesHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	messageList, findErr := self.messages.All()

	if findErr != nil {
		fmt.Printf("Error loading messages: %s\n", findErr.Error())
		encoder.Encode(messagesResponse{
			Messages: []models.Message{},
		})
		return
	}
	encoder.Encode(messagesResponse{
		Messages: messageList,
	})
}
