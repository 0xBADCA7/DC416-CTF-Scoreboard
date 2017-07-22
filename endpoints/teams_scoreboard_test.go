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

// Closed over by the TeamModelMock created by NewInMemoryTeamModel.
type teamState struct {
	Teams []models.Team
}

// Constructs a TeamModelMock that operates on an in-memory array of teams
// instead of talking to a database.
//
// Note that this mock has the quirk that, after a team with an id N gets
// deleted, all teams with ids m >= n will be decremented to id = m - 1.
func NewInMemoryTeamModel(state teamState) mocks.TeamModelMock {

	find := func(token string) (models.Team, error) {
		for _, team := range state.Teams {
			if team.SubmitToken == token {
				return team, nil
			}
		}
		team := models.Team{}
		return team, errors.New("Team not found")
	}

	all := func() ([]models.Team, error) {
		return state.Teams, nil
	}

	save := func(team *models.Team) error {
		team.Id = len(state.Teams) + 1
		state.Teams = append(state.Teams, *team)
		return nil
	}

	update := func(team *models.Team) error {
		state.Teams[team.Id-1] = *team
		return nil
	}

	del := func(team *models.Team) error {
		state.Teams = append(state.Teams[:team.Id-1], state.Teams[team.Id:]...)
		// Wouldn't happen in a DB, but we will do it to make our mock simple.
		for i := team.Id - 1; i < len(state.Teams); i++ {
			state.Teams[i].Id -= 1
		}
		return nil
	}

	return mocks.NewTeamModelMock(find, all, save, update, del)
}

func TestScoreboardEndpoint(test *testing.T) {
	state := teamState{make([]models.Team, 0)}
	teamsModel := NewInMemoryTeamModel(state)
	handler := NewTeamsScoreboardHandler(teamsModel)
	server := httptest.NewServer(handler)
	defer server.Close()

	type testFn func(TeamScoreboardResponse) error

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
		data := TeamScoreboardResponse{}
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
		return func(data TeamScoreboardResponse) error {
			numTeams := len(data.Teams)
			if numTeams != expectedLen {
				return errors.New(fmt.Sprintf("Expected %d teams. Found %d\n", expectedLen, numTeams))
			}
			return nil
		}
	}

	// Produces a function that can be passed to testOutput that checks that the Nth team
	// returned has the same fields as a given team.
	compareNthTeam := func(index int, team models.Team) testFn {
		return func(data TeamScoreboardResponse) error {
			teamFound := data.Teams[index]
			teamIsExpected := teamFound.Name == team.Name &&
				teamFound.Score == team.Score &&
				teamFound.Members == team.Members
			if !teamIsExpected {
				return errors.New(fmt.Sprintf("Expected %v to equal %v\n", teamFound, team))
			}
			return nil
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
	err = testOutput(compareNthTeam(0, team))
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
	err = testOutput(compareNthTeam(1, team2))
	if err != nil {
		test.Error(err)
	}

	// Check that deleting a team causes it to not be returned.
	teamsModel.Delete(&team)
	err = testOutput(testLenTeams(1))
	if err != nil {
		test.Error(err)
	}
	err = testOutput(compareNthTeam(0, team2))
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
	err = testOutput(compareNthTeam(0, team2))
	if err != nil {
		test.Error(err)
	}
}
