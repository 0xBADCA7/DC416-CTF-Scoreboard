package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type LoginHandler struct {
	sessions models.SessionModel
}

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Error      *string `json:"error"`
	Session    string  `json:"session"`
	RedirectTo string  `json:"redirect"`
}

func NewLoginHandler(sessions models.SessionModel) LoginHandler {
	return LoginHandler{
		sessions,
	}
}

func (self LoginHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	request := LoginRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		errMsg := "Invalid request data"
		encoder.Encode(LoginResponse{
			Error:      &errMsg,
			Session:    "",
			RedirectTo: "",
		})
		return
	}
	authErr := authentication.AdminLogin(request.Password)
	if authErr != nil {
		res.WriteHeader(http.StatusUnauthorized)
		errMsg := "Incorrect password"
		encoder.Encode(LoginResponse{
			Error:      &errMsg,
			Session:    "",
			RedirectTo: "",
		})
		return
	}
	session := models.NewSession()
	saveErr := self.sessions.Save(&session)
	if saveErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Admin login failed due to internal server error: %s\n", saveErr.Error())
		errMsg := "Server encountered an error"
		encoder.Encode(LoginResponse{
			Error:      &errMsg,
			Session:    "",
			RedirectTo: "",
		})
		return
	}
	fmt.Println("Successful admin login by", req.RemoteAddr)
	encoder.Encode(LoginResponse{
		Error:      nil,
		Session:    session.Token,
		RedirectTo: "/admin",
	})
}
