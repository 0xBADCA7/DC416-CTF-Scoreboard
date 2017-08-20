package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type AdminTeamsHandler struct {
	cfg         config.Config
	submissions models.SubmissionModel
	teams       models.TeamModel
	sessions    models.SessionModel
}

type AdminTeamsRequest struct {
	SessionToken string `json:"session"`
}

type AdminTeamsResponse struct {
	Error *string         `json:"error"`
	Teams []adminTeamInfo `json:"teams"`
}

type adminTeamInfo struct {
	Id             int      `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	SubmitToken    string   `json:"submitToken,omitempty"`
	SubmittedFlags []string `json:submittedFlags,omitempty"`
}

func NewAdminTeamsHandler(
	cfg config.Config,
	submissions models.SubmissionModel,
	teams models.TeamModel,
	sessions models.SessionModel,
) AdminTeamsHandler {
	return AdminTeamsHandler{
		cfg,
		submissions,
		teams,
		sessions,
	}
}

func errResponse(errMsg *string) AdminTeamsResponse {
	return AdminTeamsResponse{
		errMsg,
		[]adminTeamInfo{},
	}
}

func (self AdminTeamsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	request := AdminTeamsRequest{}
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := "Invalid request format."
		encoder.Encode(errResponse(&errMsg))
		return
	}
	authErr := authentication.CheckSessionToken(request.SessionToken, self.sessions)
	if authErr != nil {
		res.WriteHeader(http.StatusForbidden)
		errMsg := "Not authenticated."
		encoder.Encode(errResponse(&errMsg))
		return
	}
	teams, err := self.loadTeamInfo()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		errMsg := "Error loading team information."
		fmt.Printf("Internal error loading teams: %s\n", err.Error())
		encoder.Encode(errResponse(&errMsg))
		return
	}
	encoder.Encode(AdminTeamsResponse{
		nil,
		teams,
	})
}

// Load information about teams that we want to display on the admin page.
func (self AdminTeamsHandler) loadTeamInfo() ([]adminTeamInfo, error) {
	teamInfo := []adminTeamInfo{}
	teamList, err := self.teams.All()
	if err != nil {
		return teamInfo, err
	}
	index := 0
	// Collect information about teams and the flags they've submitted.
	for _, team := range teamList {
		submissions, err := self.submissions.All(team.Id)
		if err != nil {
			return []adminTeamInfo{}, err
		}
		teamInfo = append(teamInfo, adminTeamInfo{
			Id:             team.Id,
			Name:           team.Name,
			SubmitToken:    team.SubmitToken,
			SubmittedFlags: []string{},
		})
		// I know a loop like this is going to look inefficient, but since the number
		// of teams and flags are going to be small (surely less than 100 each),
		// this kind of clear approach is better than taking a more efficient approach.
		for _, flag := range self.cfg.Flags {
			found := false
			for _, submission := range submissions {
				if flag.Id == submission.Flag {
					teamInfo[index].SubmittedFlags = append(teamInfo[index].SubmittedFlags, flagFound)
					found = true
					break
				}
			}
			if !found {
				teamInfo[index].SubmittedFlags = append(teamInfo[index].SubmittedFlags, flagNotFound)
			}
		}
		index += 1
	}
	return teamInfo, nil
}
