package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// DeleteTeamHandler handles requests to have teams deleted by administrators.
type DeleteTeamHandler struct {
	teams    models.TeamModel
	sessions models.SessionModel
}

// NewDeleteTeamHandler constructs a DeleteTeamHandler with capabilities for managing teams and admin sessions.
func NewDeleteTeamHandler(teams models.TeamModel, sessions models.SessionModel) DeleteTeamHandler {
	return DeleteTeamHandler{
		teams,
		sessions,
	}
}

func (self DeleteTeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(models.SessionCookieName)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "You are not authorized to do that.")
		return
	}
	authErr := authentication.CheckSessionToken(cookie.Value, self.sessions)
	if authErr != nil {
		jsonError(w, http.StatusUnauthorized, "You are not authorized to do that.")
		return
	}
	ids, found := r.URL.Query()["id"]
	if !found || len(ids) == 0 {
		jsonError(w, http.StatusBadRequest, "No team specified.")
		return
	}
	id, parseErr := strconv.Atoi(ids[0])
	if parseErr != nil {
		jsonError(w, http.StatusBadRequest, "Invalid team ID.")
		return
	}
	var team models.Team
	teams, findErr := self.teams.All()
	for _, t := range teams {
		if t.Id == id {
			team = t
			break
		}
	}
	if findErr != nil || team.Id == 0 {
		jsonError(w, http.StatusBadRequest, "No team found.")
		return
	}
	deleteErr := self.teams.Delete(&team)
	if deleteErr != nil {
		jsonError(w, http.StatusInternalServerError, "Failed to delete team.")
		return
	}
	w.Write([]byte("{\"error\": null}"))
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", msg)))
}
