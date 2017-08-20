package endpoints

import (
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type MessagesPostHandler struct {
	messages models.MessageModel
	sessions models.SessionModel
}

func NewMessagesPostHandler(messages models.MessageModel, sessions models.SessionModel) MessagesPostHandler {
	return MessagesPostHandler{
		messages,
		sessions,
	}
}

func (self MessagesPostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(models.SessionCookieName)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("You are not authorized to do that."))
		return
	}
	authErr := authentication.CheckSessionToken(cookie.Value, self.sessions)
	if authErr != nil {
		err := req.ParseForm()
		res.Header().Set("Content-Type", "text/plain")
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Invalid form data submitted."))
			return
		}
		msgs, found := req.Form["message"]
		if !found || len(msgs) == 0 || len(msgs[0]) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("No message provided."))
			return
		}
		message := models.NewMessage(msgs[0])
		err = self.messages.Save(&message)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Could not save message."))
			return
		}
		res.Write([]byte("Successfully saved new message."))
	}
}
