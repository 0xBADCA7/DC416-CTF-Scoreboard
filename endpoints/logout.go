package endpoints

import (
	"database/sql"
	"net/http"
	"time"

	"../config"
	"../models"
)

// Logout instructs the browser to delete the cookie containing the
// admin's session token only if the request contains a cookie. This
// is to prevent leaking cookie information to regular users.
func Logout(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
