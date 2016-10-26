package endpoints

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"../config"
	"../teams"
)

// Register presents a page which users can use to register.
func Register(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) == "POST" {
			registerNewTeam(db, cfg, w, r)
		} else {
			registerPage(db, cfg, w, r)
		}
	}
}

// registerNewTeam handles a POST request that contains new team data and creates a team
// in the database.
func registerNewTeam(db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	team := teams.Team{}
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
		w.Write([]byte("You did submit a team name"))
		return
	}
	members, found := r.Form["members"]
	if !found || len(members) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You did submit your members' names"))
		return
	}
	team.Name = names[0]
	team.Members = members[0]
	err = team.Save(db)
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
func registerPage(db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {
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
