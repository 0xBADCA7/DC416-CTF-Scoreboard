package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type MessagesHandler struct {
	messages models.MessageModel
}

type MessagesResponse struct {
	Messages []models.Message `json:"messages"`
}

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
		encoder.Encode(MessagesResponse{
			Messages: []models.Message{},
		})
		return
	}
	encoder.Encode(MessagesResponse{
		Messages: messageList,
	})
}
