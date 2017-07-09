package endpoints

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

const adminURL = "/admin"
const flagFound = "✓ "
const flagNotFound = " "

// Contains information we will display in the table on the admin page.
// Name is the team's name, Token is the team's submission token.
// Submitted is an array of strings the same length as the number of
// flags that the scoreboard is configured with, and each index contains
// a check mark symbol if the flag was found by the team in question or
// else a space.
type adminTeamInfo struct {
	Id        int
	Name      string
	Token     string
	Submitted []string
}

// Load information about teams that we want to display on the admin page.
func loadTeamInfo(db *sql.DB, cfg *config.Config) ([]adminTeamInfo, error) {
	teamInfo := []adminTeamInfo{}
	teams, err := models.FindTeams(db)
	if err != nil {
		return teamInfo, err
	}
	index := 0
	// Collect information about teams and the flags they've submitted.
	for _, team := range teams {
		submissions, err := models.FindAllSubmissions(db, team.Id)
		if err != nil {
			return []adminTeamInfo{}, err
		}
		teamInfo = append(teamInfo, adminTeamInfo{
			Id:        team.Id,
			Name:      team.Name,
			Token:     team.SubmitToken,
			Submitted: []string{},
		})
		// I know a loop like this is going to look inefficient, but since the number
		// of teams and flags are going to be small (surely less than 100 each),
		// this kind of clear approach is better than taking a more efficient approach.
		for _, flag := range cfg.Flags {
			found := false
			for _, submission := range submissions {
				if flag.Id == submission.Flag {
					teamInfo[index].Submitted = append(teamInfo[index].Submitted, flagFound)
					found = true
					break
				}
			}
			if !found {
				teamInfo[index].Submitted = append(teamInfo[index].Submitted, flagNotFound)
			}
		}
		index += 1
	}
	return teamInfo, nil
}

// Admin handles requests for /admin, redirecting to / if no session token
// is received before testing whether the user has authenticated.
// Authenticated users are served a page containing secret information
// about teams.
func Admin(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authErr := authentication.CheckSessionToken(r, db)
		if authErr != nil {
			http.SetCookie(w, &http.Cookie{
				Name:    models.SessionCookieName,
				Value:   "",
				Expires: time.Now(),
			})
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("Access Denied"))
			return
		}
		t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "admin.html"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("content-Type", "text/plain")
			w.Write([]byte("Could not load admin page"))
			return
		}
		teams, err := loadTeamInfo(db, cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("content-Type", "text/plain")
			fmt.Println("Error loading admin data on teams:", err.Error())
			w.Write([]byte("Could not load admin data"))
			return
		}
		data := struct {
			Teams []adminTeamInfo
			Flags []config.Flag
		}{
			teams,
			cfg.Flags,
		}
		t.Execute(w, data)
	}
}
