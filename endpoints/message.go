package endpoints

import (
	"database/sql"
	"net/http"
	"strings"

	"../authentication"
	"../config"
	"../models"
)

// PostMessage handles requests to create new messages.
func PostMessage(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authErr := authentication.CheckSessionToken(r, db)
		if authErr != nil || strings.ToUpper(r.Method) != "POST" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			saveMessage(db, cfg, w, r)
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
