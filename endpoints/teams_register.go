package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// RegisterTeamHandler handles requests to register new teams.
type RegisterTeamHandler struct {
	teams    models.TeamModel
	sessions models.SessionModel
}

type registerTeamRequest struct {
	SessionToken string `json:"session"`
	TeamName     string `json:"name"`
	TeamMembers  string `json:"members"`
}

type registerTeamResponse struct {
	Error       *string `json:"error"`
	SubmitToken string  `json:"submitToken"`
}

// NewRegisterTeamHandler constructs a RegisterTeamHandler with the capabilities to manage teams and access
// information about administrator sessions.
func NewRegisterTeamHandler(teams models.TeamModel, sessions models.SessionModel) RegisterTeamHandler {
	return RegisterTeamHandler{
		teams,
		sessions,
	}
}

func (self RegisterTeamHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	request := registerTeamRequest{}
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := "Invalid request."
		encoder.Encode(registerTeamResponse{
			&errMsg,
			"",
		})
		return
	}
	authErr := authentication.CheckSessionToken(request.SessionToken, self.sessions)
	if authErr != nil {
		res.WriteHeader(http.StatusUnauthorized)
		errMsg := "You are not authorized to do that."
		encoder.Encode(registerTeamResponse{
			&errMsg,
			"",
		})
		return
	}
	team := models.Team{}
	team.Name = request.TeamName
	team.Members = request.TeamMembers
	err := self.teams.Save(&team)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error saving team: %s\n", err.Error())
		errMsg := "Error saving team."
		encoder.Encode(registerTeamResponse{
			&errMsg,
			"",
		})
		return
	}
	encoder.Encode(registerTeamResponse{
		nil,
		team.SubmitToken,
	})
}
