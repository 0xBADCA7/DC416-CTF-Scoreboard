package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// TeamScoreboardHandler handles GET requests made to retrieve information about teams
// that can be displayed on a scoreboard.
type TeamScoreboardHandler struct {
	teams models.TeamModel
}

// teamScoreboardResponse contains the information written to the client.
type teamScoreboardResponse struct {
	Error *string      `json:"error"`
	Teams []SBTeamInfo `json:"teams"`
}

// SBTeamInfo contains information about a team to display on a scoreboard.
type SBTeamInfo struct {
	Name           string `json:"name"`
	Score          int    `json:"score"`
	Position       int    `json:"position"`
	Members        string `json:"members"`
	LastSubmission string `json:"lastSubmission"`
}

// NewTeamsScoreboardHandler constructs a new TeamScoreboardHandler.
func NewTeamsScoreboardHandler(teams models.TeamModel) TeamScoreboardHandler {
	return TeamScoreboardHandler{
		teams,
	}
}

func (self TeamScoreboardHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(res)
	teams, err := self.teams.All()
	if err != nil {
		fmt.Printf("Error retrieving teams for scoreboard: %v\n", err)
		errMsg := "Could not retrieve information about teams"
		encoder.Encode(teamScoreboardResponse{
			Error: &errMsg,
			Teams: []SBTeamInfo{},
		})
		return
	}
	sort.Sort(models.TeamByScore(teams))
	teamInfo := teamScoreboardResponse{
		Error: nil,
		Teams: make([]SBTeamInfo, len(teams)),
	}
	for i, team := range teams {
		lastSubmitted := "No flags submitted yet"
		if team.LastSubmission.After(time.Date(2015, time.January, 1, 1, 0, 0, 0, time.UTC)) {
			lastSubmitted = team.LastSubmission.Local().Format(time.UnixDate)
		}
		teamInfo.Teams[i] = SBTeamInfo{
			Name:           team.Name,
			Score:          team.Score,
			Position:       i + 1,
			Members:        team.Members,
			LastSubmission: lastSubmitted,
		}
	}
	encoder.Encode(teamInfo)
}
