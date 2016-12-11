package endpoints

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"../config"
)

// Login presents routes GET requests to serve a login page for admins and POST
// requests to handle a login form submission.
func Login(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) == "POST" {
			adminLogin(db, cfg, w, r)
		} else {
			loginPage(db, cfg, w, r)
		}
	}
}

// adminLogin checks the password provided to the application against a configured password,
// granting access to the admin dashboard if the credentials are correct.
func adminLogin(db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	// TODO - Implement me
	fmt.Println("Ignoring login submission")
	w.Write([]byte("Rejected"))
}

// loginPage serves a page containing a login form for admins to access the admin dashboard.
func loginPage(db *sql.DB, cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got a request for the admin login page")
	t, err := template.ParseFiles(path.Join(cfg.HTMLDir, "login.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Could not load login page"))
		return
	}
	err = t.Execute(w, nil)
}
