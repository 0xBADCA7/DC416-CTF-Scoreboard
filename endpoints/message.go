package endpoints

import (
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type MessageHandler struct {
	cfg      config.Config
	messages models.MessageModel
	sessions models.SessionModel
}

func NewMessageHandler(cfg config.Config, messages models.MessageModel, sessions models.SessionModel) MessageHandler {
	return MessageHandler{
		cfg,
		messages,
		sessions,
	}
}

func (self MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authErr := authentication.CheckSessionToken(r, self.sessions)
	if authErr == nil && strings.ToUpper(r.Method) == "POST" {
		saveMessage(&self.cfg, self.messages, w, r)
	} else {
		messagesPage(&self.cfg, self.messages, w, r)
	}
}

func saveMessage(cfg *config.Config, messages models.MessageModel, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid form data submitted."))
		return
	}
	msgs, found := r.Form["message"]
	if !found || len(msgs) == 0 || len(msgs[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No message provided."))
		return
	}
	message := models.NewMessage(msgs[0])
	err = messages.Save(&message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not save message."))
		return
	}
	w.Write([]byte("Successfully saved new message."))
}

func messagesPage(cfg *config.Config, messages models.MessageModel, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "messages.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Could not load messages."))
		return
	}
	messageList, findErr := messages.All()
	if findErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Could not load messages."))
		return
	}
	data := struct {
		Messages []models.Message
	}{messageList}
	t.Execute(w, data)
}
