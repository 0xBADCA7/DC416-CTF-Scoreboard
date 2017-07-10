package endpoints

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// RegistrationHandler implements http.ServeHTTP to handle GET requests, which it responds to with
// a page containing a registration form, and POST requests, which it handles by registering a new team.
type RegistrationHandler struct {
	cfg      config.Config
	teams    models.TeamModel
	sessions models.SessionModel
}

// NewRegistrationHandler construcs a new RegistrationHandler with a means of managing teams and sessions.
func NewRegistrationHandler(cfg config.Config, teams models.TeamModel, sessions models.SessionModel) RegistrationHandler {
	return RegistrationHandler{
		cfg,
		teams,
		sessions,
	}
}

// ServeHTTP handles requests to either view a registration form or upload a new team's info.
func (self RegistrationHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	authErr := authentication.CheckSessionToken(req, self.sessions)
	if authErr != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	} else if strings.ToUpper(req.Method) == "POST" {
		registerNewTeam(&self.cfg, self.teams, res, req)
	} else {
		registerPage(&self.cfg, res, req)
	}
}

// registerNewTeam handles a POST request that contains new team data and creates a team
// in the database.
func registerNewTeam(cfg *config.Config, teams models.TeamModel, w http.ResponseWriter, r *http.Request) {
	team := models.Team{}
	fmt.Println("Got a POST request to register a new team")
	err := r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You submitted invalid form data"))
		return
	}
	fmt.Println("Got POST data", r.Form)
	names, found := r.Form["name"]
	if !found || len(names) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You did not submit a team name"))
		return
	}
	members, found := r.Form["members"]
	if !found || len(members) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You did not submit your members' names"))
		return
	}
	team.Name = names[0]
	team.Members = members[0]
	err = teams.Save(&team)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not create your team. Please make sure your name is not taken."))
		return
	}
	msg := fmt.Sprintf(`Your team has successfully been registered.
Your submission token is %s.
Please make sure not to lose or share it with anyone not on your team.`,
		team.SubmitToken)
	w.Write([]byte(msg))
}

// registerPage serves the register.html page which contains a form that users can fill out
// to register their team.
func registerPage(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got a request for the register page")
	t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "register.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Could not load register page"))
		return
	}
	err = t.Execute(w, nil)
}
