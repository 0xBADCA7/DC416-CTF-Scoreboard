package endpoints

import (
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type LoginHandler struct {
	sessions models.SessionModel
}

func NewLoginHandler(sessions models.SessionModel) LoginHandler {
	return LoginHandler{
		sessions,
	}
}

func (self LoginHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	badPwdMsg := []byte("Incorrect password")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Set("Content-Type", "text/plain")
		res.Write(badPwdMsg)
		return
	}
	password, found := req.Form["password"]
	if !found {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Set("Content-Type", "text/plain")
		res.Write(badPwdMsg)
		return
	}
	authErr := authentication.AdminLogin(password[0])
	if authErr != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Set("Content-Type", "text/plain")
		res.Write(badPwdMsg)
		return
	}
	session := models.NewSession()
	saveErr := self.sessions.Save(&session)
	if saveErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte(fmt.Sprintf("Could not log in. Reason: %s\n", saveErr.Error())))
		return
	}
	fmt.Println("Successful admin login by", req.RemoteAddr)
	http.SetCookie(res, &http.Cookie{
		Name:    models.SessionCookieName,
		Value:   session.Token,
		Expires: session.Expires,
	})
	// adminURL defined in endpoints/admin.go
	http.Redirect(res, req, adminURL, http.StatusSeeOther)
}
