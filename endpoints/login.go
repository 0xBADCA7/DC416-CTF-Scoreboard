package endpoints

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
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
	err := r.ParseForm()
	badPwdMsg := []byte("Incorrect password")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(badPwdMsg)
		return
	}
	password, found := r.Form["password"]
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(badPwdMsg)
		return
	}
	authErr := authentication.AdminLogin(db, password[0])
	if authErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(badPwdMsg)
		return
	}
	session := models.NewSession()
	saveErr := session.Save(db)
	if saveErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("Could not log in. Reason: %s\n", saveErr.Error())))
		return
	}
	fmt.Println("Successful admin login by", r.RemoteAddr)
	http.SetCookie(w, &http.Cookie{
		Name:    models.SessionCookieName,
		Value:   session.Token,
		Expires: session.Expires,
	})
	// adminURL defined in endpoints/admin.go
	http.Redirect(w, r, adminURL, http.StatusSeeOther)
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
