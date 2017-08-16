package endpoints

import (
	"html/template"
	"net/http"
	"path"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
)

type LoginPageHandler struct {
	cfg config.Config
}

func NewLoginPageHandler(cfg config.Config) LoginPageHandler {
	return LoginPageHandler{
		cfg,
	}
}

func (self LoginPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(path.Join(self.cfg.HTMLDir, "login.html"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("Could not load login page"))
		return
	}
	err = t.Execute(res, nil)
}
