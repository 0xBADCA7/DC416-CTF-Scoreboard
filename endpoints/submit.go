package endpoints

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// SubmissionHandler implements http.ServeHTTP to handle GET requests, which it responds to with a page listing
// a submit form or POST requests, which it handles by checking a flag.
type SubmissionHandler struct {
	cfg         config.Config
	submissions models.SubmissionModel
	teams       models.TeamModel
}

// NewSubmissionHandler constructs a new submission handler with a means of managing submissions and teams..
func NewSubmissionHandler(cfg config.Config, subs models.SubmissionModel, teams models.TeamModel) SubmissionHandler {
	return SubmissionHandler{
		cfg,
		subs,
		teams,
	}
}

// ServeHTTP handles requests to either view a submission form or upload a new flag.
func (self SubmissionHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if strings.ToUpper(req.Method) == "POST" {
		handleSubmission(&self.cfg, self.submissions, self.teams, res, req)
	} else {
		submitPage(&self.cfg, res, req)
	}
}

// submitPage serves the HTML page allowing users to submit flags.
func submitPage(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "submit.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Could not load register page"))
		return
	}
	err = t.Execute(w, nil)
}

// handleSubmission handles POST requests to /submit, issued by users when they are trying to submit
// a flag. It prevents teams from entering the same flag multiple times and makes sure that the
// submission token submitted is valid.
func handleSubmission(cfg *config.Config, submissionModel models.SubmissionModel, teamModel models.TeamModel, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got a request to submit a flag")
	w.Header().Set("Content-Type", "text/plain")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Your input is poorly formatted!"))
		return
	}
	fmt.Println(r.Form)
	tokens, found := r.Form["token"]
	if !found || len(tokens) == 0 {
		fmt.Println("Missing token")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing the token field. Please supply the submission token you were assigned."))
		return
	}
	flags, found := r.Form["flag"]
	if !found || len(flags) == 0 {
		fmt.Println("Missing flag")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing the flag field. Please supply secret flag."))
		return
	}
	team, err := teamModel.Find(tokens[0])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You submitted an invalid token. Please make sure you entered it correctly."))
		return
	}
	flag := config.Flag{}
	found = false
	for _, f := range cfg.Flags {
		if f.Secret == flags[0] {
			flag = f
			found = true
		}
	}
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The flag you submitted is invalid. Please check that it is formatted correctly."))
		return
	}
	submission, err := submissionModel.Find(team.Id, flag.Id)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You cannot submit the same flag multiple times."))
		return
	}
	submission.Flag = flag.Id
	submission.Owner = team.Id
	err = submissionModel.Save(&submission)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not record your submission. Please notify the CTF administrators."))
		return
	}
	team.Score += flag.Reward
	team.LastSubmission = time.Now()
	err = teamModel.Update(&team)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not update your score. Please notify the CTF administrators."))
		return
	}
	w.Write([]byte(fmt.Sprintf(
		"Congrats! You have been awarded %d points. Your score is now %d.\n",
		flag.Reward,
		team.Score)))

}
