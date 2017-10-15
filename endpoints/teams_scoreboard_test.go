package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/mocks"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

func TestScoreboardEndpoint(test *testing.T) {
	teamsModel := mocks.NewInMemoryTeamModel()
	handler := NewTeamsScoreboardHandler(teamsModel)
	server := httptest.NewServer(handler)
	defer server.Close()

	type testFn func(teamScoreboardResponse) error

	// Makes a request that gets handled the team scoreboard handler and runs a test function that
	// can inspect (upon successful decoding) the JSON returned by the handler.
	testOutput := func(runTest testFn) error {
		response, err := http.Get(server.URL)
		if err != nil {
			return err
		}
		if response.StatusCode != 200 {
			return errors.New(fmt.Sprintf("Expected status code 200. Got %d\n", response.StatusCode))
		}
		data := teamScoreboardResponse{}
		decoder := json.NewDecoder(response.Body)
		defer response.Body.Close()
		err = decoder.Decode(&data)
		if err != nil {
			return err
		}
		return runTest(data)
	}

	// Produces a function that can be passed to testOutput to check the number of teams returned.
	testLenTeams := func(expectedLen int) testFn {
		return func(data teamScoreboardResponse) error {
			numTeams := len(data.Teams)
			if numTeams != expectedLen {
				return errors.New(fmt.Sprintf("Expected %d teams. Found %d\n", expectedLen, numTeams))
			}
			return nil
		}
	}

	// Produces a function that can be passed to testOutput that checks if a team is presnt
	// in the response from the server.
	tryFindTeam := func(team models.Team) testFn {
		return func(data teamScoreboardResponse) error {
			for _, teamFound := range data.Teams {
				teamIsExpected := teamFound.Name == team.Name &&
					teamFound.Score == team.Score &&
					teamFound.Members == team.Members
				if teamIsExpected {
					return nil
				}
			}
			return errors.New(fmt.Sprintf("Expected %v to be in response\n", team))
		}
	}

	// Check that we start with no teams
	err := testOutput(testLenTeams(0))
	if err != nil {
		test.Error(err)
	}

	// Check that we get the score of the first team that is registered.
	// Also effectively tests idempotence.
	team := models.Team{
		Id:             0,
		Name:           "first",
		Members:        "",
		Score:          0,
		SubmitToken:    "testtoken",
		LastSubmission: time.Now(),
	}
	teamsModel.Save(&team)
	err = testOutput(testLenTeams(1))
	if err != nil {
		test.Error(err)
	}
	err = testOutput(tryFindTeam(team))
	if err != nil {
		test.Error(err)
	}

	// Check that inserting another team causes two to be returned.
	team2 := models.Team{
		Id:             0,
		Name:           "second",
		Members:        "",
		Score:          30,
		SubmitToken:    "testtoken2",
		LastSubmission: time.Now(),
	}
	teamsModel.Save(&team2)
	err = testOutput(testLenTeams(2))
	if err != nil {
		test.Error(err)
	}
	err = testOutput(tryFindTeam(team2))
	if err != nil {
		test.Error(err)
	}

	// Check that deleting a team causes it to not be returned.
	teamsModel.Delete(&team)
	err = testOutput(testLenTeams(1))
	if err != nil {
		test.Error(err)
	}
	err = testOutput(tryFindTeam(team2))
	if err != nil {
		test.Error(err)
	}

	// Check that update a team causes its updated info to be returned.
	team2.Score += 30
	teamsModel.Update(&team2)
	err = testOutput(testLenTeams(1))
	if err != nil {
		test.Error(err)
	}
	err = testOutput(tryFindTeam(team2))
	if err != nil {
		test.Error(err)
	}
}
