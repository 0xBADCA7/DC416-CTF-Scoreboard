package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type TeamSubmitHandler struct {
	teams       models.TeamModel
	submissions models.SubmissionModel
	flags       []config.Flag
}

type TeamSubmitRequest struct {
	Token string `json:"token"`
	Flag  string `json:"flag"`
}

type TeamSubmitResponse struct {
	Error     *string `json:"error"`
	IsCorrect bool    `json:"correct"`
	NewScore  int     `json:"newScore"`
}

func NewTeamSubmitHandler(teams models.TeamModel, submissions models.SubmissionModel, flags []config.Flag) TeamSubmitHandler {
	return TeamSubmitHandler{
		teams,
		submissions,
		flags,
	}
}

func (self TeamSubmitHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	decoder := json.NewDecoder(req.Body)
	request := TeamSubmitRequest{}
	defer req.Body.Close()
	err := decoder.Decode(&request)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprintf("Error submitting flag: %v\n", err)
		encoder.Encode(TeamSubmitResponse{
			Error:     &errMsg,
			IsCorrect: false,
			NewScore:  0,
		})
		return
	}
	team, err := self.teams.Find(request.Token)
	if err != nil {
		res.WriteHeader(http.StatusForbidden)
		errMsg := "Invalid submission token"
		encoder.Encode(TeamSubmitResponse{
			Error:     &errMsg,
			IsCorrect: false,
			NewScore:  0,
		})
		return
	}
	found := false
	flag := config.Flag{}
	for _, _flag := range self.flags {
		if _flag.Secret == request.Flag {
			flag = _flag
			found = true
			break
		}
	}
	if !found {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := "The flag you submitted is not correct."
		encoder.Encode(TeamSubmitResponse{
			Error:     &errMsg,
			IsCorrect: false,
			NewScore:  team.Score,
		})
		return
	}
	// IMPORTANT
	// I decided to try to save a record of the team's submission before updating their score here because,
	// if we somehow encountered a bug allowing a user to trigger a database operation failure just before we
	// save their submission, they could exploit this to get as many points as they want.  This way, we can
	// resolve failures to update a user's score manually and with no risk.
	submission := models.NewSubmission(flag.Id, team.Id)
	saveErr := self.submissions.Save(&submission)
	if saveErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := "Failed to save submission."
		encoder.Encode(TeamSubmitResponse{
			Error:     &errMsg,
			IsCorrect: true,
			NewScore:  team.Score,
		})
		return
	}
	team.Score += flag.Reward
	team.LastSubmission = time.Now()
	updateErr := self.teams.Update(&team)
	if updateErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		errMsg := "Internal error; failed to update your score. Please contact an administrator."
		encoder.Encode(TeamSubmitResponse{
			Error:     &errMsg,
			IsCorrect: true,
			NewScore:  team.Score - flag.Reward,
		})
		return
	}
	encoder.Encode(TeamSubmitResponse{
		Error:     nil,
		IsCorrect: true,
		NewScore:  team.Score,
	})
}
