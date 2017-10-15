package endpoints

import (
	"html/template"
	"net/http"
	"path"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// RegisterPageHandler handles requests to retrieve the page on which teams can be registered.
type RegisterPageHandler struct {
	cfg      config.Config
	sessions models.SessionModel
}

// NewRegisterPageHandler constructs a new RegisterPageHandler with the capability to access information about
// administrator sessions.
func NewRegisterPageHandler(cfg config.Config, sessions models.SessionModel) RegisterPageHandler {
	return RegisterPageHandler{
		cfg,
		sessions,
	}
}

func (self RegisterPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(models.SessionCookieName)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("You are not allowed to do that."))
		return
	}
	authErr := authentication.CheckSessionToken(cookie.Value, self.sessions)
	if authErr != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("You are not allowed to do that."))
		return
	}
	t, err := template.ParseFiles(path.Join(self.cfg.HTMLDir, "register.html"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("Could not load register page"))
		return
	}
	err = t.Execute(res, nil)
}
