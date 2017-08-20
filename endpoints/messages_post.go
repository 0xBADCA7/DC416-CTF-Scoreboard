package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type MessagesPostHandler struct {
	messages models.MessageModel
	sessions models.SessionModel
}

type MessagesPostRequest struct {
	SessionToken   string `json:"session"`
	MessageContent string `json:"content"`
}

type MessagesPostResponse struct {
	Error *string `json:"error"`
}

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

	request := MessagesPostRequest{}
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := "Invalid request."
		encoder.Encode(MessagesPostResponse{
			&errMsg,
		})
		return
	}
	authErr := authentication.CheckSessionToken(request.SessionToken, self.sessions)
	if authErr != nil {
		res.WriteHeader(http.StatusUnauthorized)
		errMsg := "You are not authorized to do that."
		encoder.Encode(MessagesPostResponse{
			&errMsg,
		})
		return
	}
	message := models.NewMessage(request.MessageContent)
	err := self.messages.Save(&message)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Error saving new message: %s\n", err.Error())
		encoder.Encode(MessagesPostResponse{
			&errMsg,
		})
		return
	}
	encoder.Encode(MessagesPostResponse{
		nil,
	})
}
