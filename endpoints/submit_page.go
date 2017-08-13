package endpoints

import (
	"html/template"
	"net/http"
	"path"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
)

// SubmitPageHandler implements http.ServeHTTP to handle GET requests, which it responds to with a page listing
// a submit form or POST requests, which it handles by checking a flag.
type SubmitPageHandler struct {
	cfg config.Config
}

// NewSubmitPageHandler constructs a new submission handler with a means of managing submissions and teams..
func NewSubmitPageHandler(cfg config.Config) SubmitPageHandler {
	return SubmitPageHandler{
		cfg,
	}
}

// ServeHTTP handles requests to either view a submission form or upload a new flag.
func (self SubmitPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(path.Join(self.cfg.HTMLDir, "submit.html"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("Could not load register page"))
		return
	}
	err = t.Execute(res, nil)
}
