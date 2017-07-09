package endpoints

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// PostMessage handles requests to create new messages.
func PostMessage(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authErr := authentication.CheckSessionToken(r, db)
		if authErr == nil && strings.ToUpper(r.Method) == "POST" {
			saveMessage(db, cfg, w, r)
		} else {
			messagesPage(db, cfg, w, r)
		}
	}
}

func saveMessage(db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {
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
	err = message.Save(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not save message."))
		return
	}
	w.Write([]byte("Successfully saved new message."))
}

func messagesPage(db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "messages.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Could not load messages."))
		return
	}
	messages, findErr := models.AllMessages(db)
	if findErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Could not load messages."))
		return
	}
	fmt.Println("got messages ", messages)
	data := struct {
		Messages []models.Message
	}{messages}
	t.Execute(w, data)
}
