package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/mocks"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

func TestSubmitEndpoint(test *testing.T) {
	teamsModel := mocks.NewInMemoryTeamModel()
	submissionsModel := mocks.NewInMemorySubmissionModel()
	handler := NewTeamSubmitHandler(teamsModel, submissionsModel, []config.Flag{
		{1, "flag{secret1}", 10},
		{2, "flag{m00nc4k3}", 30},
		{3, "flag{h3ll0w0rld}", 50},
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	team1 := models.Team{
		Id:             1,
		Name:           "first",
		Members:        "1",
		Score:          0,
		SubmitToken:    "team1token",
		LastSubmission: time.Now(),
	}
	team2 := models.Team{
		Id:             2,
		Name:           "second",
		Members:        "2",
		Score:          0,
		SubmitToken:    "team2token",
		LastSubmission: time.Now(),
	}
	teamsModel.Save(&team1)
	teamsModel.Save(&team2)

	testCases := []struct {
		ExpectedStatus         int
		ExpectedNumSubmissions int
		TeamId                 int
		SubmitToken            string
		FlagToSubmit           string
		ExpectError            bool
		ExpectFlagValid        bool
		ExpectedNewScore       int
	}{
		{200, 1, 1, "team1token", "flag{secret1}", false, true, 10},
		{200, 2, 1, "team1token", "flag{h3ll0w0rld}", false, true, 60},
		{200, 1, 2, "team2token", "flag{secret1}", false, true, 10},
		{400, 1, 2, "team2token", "flag{h3ll0world}", true, false, 10},
		{403, 0, 0, "team2t0k3n", "flag{hell0w0rld}", true, false, 0},
		{400, 2, 1, "team1token", "flag{secret1}", true, true, 60},
	}

	for i, testCase := range testCases {
		test.Logf("Running test case #%d\n", i)
		reqData := teamSubmitRequest{
			Token: testCase.SubmitToken,
			Flag:  testCase.FlagToSubmit,
		}
		encoded, _ := json.Marshal(&reqData)
		response, err := http.Post(server.URL, "application/json", bytes.NewReader(encoded))
		if err != nil {
			test.Error(err)
		}
		if response.StatusCode != testCase.ExpectedStatus {
			test.Errorf("Expected status %d. Got %d\n", testCase.ExpectedStatus, response.StatusCode)
		}
		data := teamSubmitResponse{}
		decoder := json.NewDecoder(response.Body)
		defer response.Body.Close()
		err = decoder.Decode(&data)
		if err != nil {
			test.Error(err)
		}
		gotErr := data.Error != nil
		if gotErr && !testCase.ExpectError {
			test.Errorf("Got unexpected error: '%v'\n", *data.Error)
		}
		if !gotErr && testCase.ExpectError {
			test.Errorf("Did not get error where it was expected")
		}
		if data.IsCorrect != testCase.ExpectFlagValid {
			test.Errorf("Expected flag to be valid? %v. Was valid? %v.", testCase.ExpectFlagValid, data.IsCorrect)
		}
		if data.NewScore != testCase.ExpectedNewScore {
			test.Errorf("Expected new score to be %d. Got %d.", testCase.ExpectedNewScore, data.NewScore)
		}
		submittedFlags, err := submissionsModel.All(testCase.TeamId)
		if err != nil {
			test.Error(err)
		}
		if len(submittedFlags) != testCase.ExpectedNumSubmissions {
			test.Errorf("Expected %d submissions. Found %d.\n", testCase.ExpectedNumSubmissions, len(submittedFlags))
		}
	}
}
