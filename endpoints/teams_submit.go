package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		fmt.Printf("Request handler got decode error %s\n", err.Error())
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
		fmt.Printf("Request handler got lookup error %s\n", err.Error())
		res.WriteHeader(http.StatusForbidden)
		errMsg := "Invalid submission token"
		fmt.Printf("ERROR: User with IP %s submitted invalid submit token %s\n", req.RemoteAddr, request.Token)
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
		fmt.Printf("Did not find flag %s\n", request.Flag)
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
	fmt.Printf("Saving submission for flag %s\n", request.Flag)
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
	fmt.Println("Saved successfully.")
	fmt.Println("Updating team")
	team.Score += flag.Reward
	updateErr := self.teams.Update(&team)
	if updateErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		errMsg := "Internal error; failed to update your score. Please contact an administrator."
		fmt.Printf("Failed to update score of team %d. Error: %v\n", team.Id, updateErr)
		encoder.Encode(TeamSubmitResponse{
			Error:     &errMsg,
			IsCorrect: true,
			NewScore:  team.Score - flag.Reward,
		})
		return
	}
	fmt.Println("Updated successfully")
	encoder.Encode(TeamSubmitResponse{
		Error:     nil,
		IsCorrect: true,
		NewScore:  team.Score,
	})
}
