package endpoints

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"../config"
	"../teams"
)

// Index creates a request handler that serves index.html, the main scoreboard page with all of
// the teams and their scores.
func Index(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "index.html"))
		if err != nil {
			fmt.Println("Error parsing template", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("We done goofed! Try again in a few minutes."))
			return
		}
		teamInfo, err := teams.FindTeams(db)
		if err != nil {
			fmt.Println("Error finding teams", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("We done goofed! Try again in a few minutes."))
			return
		}
		data := struct {
			Message string
			Teams   []teams.TeamScore
		}{
			"Hello DC416!",
			[]teams.TeamScore{},
		}
		fmt.Println("Got teams", teamInfo)
		for _, team := range teamInfo {
			data.Teams = append(data.Teams, teams.TeamScore{team.Name, team.Members, team.Score})
		}
		t.Execute(w, data)
	}
}
