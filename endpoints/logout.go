package endpoints

import (
	"net/http"
	"time"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type LogoutHandler struct {
	cfg      config.Config
	sessions models.SessionModel
}

func NewLogoutHandler(cfg config.Config, sessions models.SessionModel) LogoutHandler {
	return LogoutHandler{
		cfg,
		sessions,
	}
}

// ServeHTTP instructs the browser to delete the cookie containing the
// admin's session token only if the request contains a cookie. This
// is to prevent leaking cookie information to regular users.
func (self LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(models.SessionCookieName)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Only send "Set-Cookie: ;Expires: now" if a cookie was sent.
	http.SetCookie(w, &http.Cookie{
		Name:    models.SessionCookieName,
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
