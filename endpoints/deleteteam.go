package endpoints

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"../authentication"
	"../config"
	"../models"
)

// DeleteTeam handles requests to delete a team issued by an admin.
func DeleteTeam(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authErr := authentication.CheckSessionToken(r, db)
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
		teams, findErr := models.FindTeams(db)
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
		deleteErr := team.Delete(db)
		if deleteErr != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to delete team.")
			return
		}
		w.Write([]byte("{\"error\": null}"))
	}
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", msg)))
}
