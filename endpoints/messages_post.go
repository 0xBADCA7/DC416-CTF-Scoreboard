package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// MessagesPostHandler handles requests to post new messages.
type MessagesPostHandler struct {
	messages models.MessageModel
	sessions models.SessionModel
}

type messagesPostRequest struct {
	SessionToken   string `json:"session"`
	MessageContent string `json:"content"`
}

type messagesPostResponse struct {
	Error *string `json:"error"`
}

// NewMessagesPostHandler constructs a new MessagesPostHandler with capabilities for working with messages and
// administrator sessions.
func NewMessagesPostHandler(messages models.MessageModel, sessions models.SessionModel) MessagesPostHandler {
	return MessagesPostHandler{
		messages,
		sessions,
	}
}

func (self MessagesPostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(res)
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	request := messagesPostRequest{}
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := "Invalid request."
		encoder.Encode(messagesPostResponse{
			&errMsg,
		})
		return
	}
	authErr := authentication.CheckSessionToken(request.SessionToken, self.sessions)
	if authErr != nil {
		res.WriteHeader(http.StatusUnauthorized)
		errMsg := "You are not authorized to do that."
		encoder.Encode(messagesPostResponse{
			&errMsg,
		})
		return
	}
	message := models.NewMessage(request.MessageContent)
	err := self.messages.Save(&message)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Error saving new message: %s\n", err.Error())
		encoder.Encode(messagesPostResponse{
			&errMsg,
		})
		return
	}
	encoder.Encode(messagesPostResponse{
		nil,
	})
}
