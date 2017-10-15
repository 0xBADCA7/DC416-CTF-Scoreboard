package endpoints

import (
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

const adminURL = "/admin"
const flagFound = "✓ "
const flagNotFound = " "

// AdminPageHandler handles requests to
type AdminPageHandler struct {
	cfg         config.Config
	submissions models.SubmissionModel
	teams       models.TeamModel
	sessions    models.SessionModel
}

// NewAdminPageHandler constructs an AdminPageHandler capable of dealing with submissions, teams, and admin sessions.
func NewAdminPageHandler(
	cfg config.Config,
	submissions models.SubmissionModel,
	teams models.TeamModel,
	sessions models.SessionModel) AdminPageHandler {
	return AdminPageHandler{
		cfg,
		submissions,
		teams,
		sessions,
	}
}

func (self AdminPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(models.SessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Not authenticated."))
		return
	}
	authErr := authentication.CheckSessionToken(cookie.Value, self.sessions)
	if authErr != nil {
		http.SetCookie(w, &http.Cookie{
			Name:    models.SessionCookieName,
			Value:   "",
			Expires: time.Now(),
		})
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Access Denied"))
		return
	}
	t, err := template.ParseFiles(path.Join(self.cfg.HTMLDir, "admin.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-Type", "text/plain")
		w.Write([]byte("Could not load admin page"))
		return
	}
	t.Execute(w, nil)
}
