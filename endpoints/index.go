package endpoints

import (
	"database/sql"
	"html/template"
	"net/http"
	"path"

	"../config"
)

// Index creates a request handler that serves index.html, the main scoreboard page with all of
// the teams and their scores.
func Index(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "index.html"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("We done goofed! Try again in a few minutes."))
			return
		}
		data := struct {
			Message string
		}{
			"Hello, DC416!",
		}
		t.Execute(w, data)
	}
}
