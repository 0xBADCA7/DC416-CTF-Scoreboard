package authentication

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	auth "github.com/StratumSecurity/scryptauth"

	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// passwordEnvVar is the name of the environment variable that we expect
// the password for admin access to be provided under.
const passwordEnvVar = "CTF_PASSWORD"

// ErrExpiredToken is an error that indicates that a session token is expired.
var ErrExpiredToken = errors.New("submitted token is expired")

// CheckAuthorization determines whether a submitted session token is
// valid and not expired.
func CheckAuthorization(db *sql.DB, submittedToken string) error {
	session, findErr := models.FindSession(db, submittedToken)
	if findErr != nil {
		return findErr
	}
	if session.IsExpired() {
		return ErrExpiredToken
	}
	return nil
}

// CheckSessionToken looks for an appropriately named `session` cookie
// in the provided request and then tests whether the session id
// sent is valid.
func CheckSessionToken(r *http.Request, db *sql.DB) error {
	sessionCookie, err := r.Cookie(models.SessionCookieName)
	if err != nil {
		return err
	}
	return CheckAuthorization(db, sessionCookie.Value)
}

// HashAdminPassword applies a secure scrypt-based password hash
// to the value contained in the environment variable used to
// supply an admin password.  If no password is provided, the
// variable is left as the empty string.
func HashAdminPassword() {
	adminPwd := getAdminPassword()
	if len(adminPwd) > 0 {
		hashParams := auth.DefaultHashConfiguration()
		hashed, hashErr := auth.GenerateFromPassword([]byte(adminPwd), hashParams)
		if hashErr != nil {
			panic(hashErr)
		}
		os.Setenv(passwordEnvVar, string(hashed))
	}
}

// AdminLogin checks if a supplied password matches the one the
// scoreboard is configured to restrict access to the admin page with.
func AdminLogin(db *sql.DB, password string) error {
	expected := getAdminPassword()
	return auth.CompareHashAndPassword([]byte(expected), []byte(password))
}

// getAdminPassword obtains the password required to access the
// admin dashboard.  It should be hashed at the start of main using
// the HashAdminPassword function.  An empty string indicates that
// no password was set and all admin operations should be rejected.
func getAdminPassword() string {
	return os.Getenv(passwordEnvVar)
}
