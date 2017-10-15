package endpoints

import (
	"html/template"
	"net/http"
	"path"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
)

// MessagePageHandler handles requests to view the messages page.
type MessagePageHandler struct {
	cfg config.Config
}

// NewMessagePageHandler constructs a new MessagePageHandler.
func NewMessagePageHandler(cfg config.Config) MessagePageHandler {
	return MessagePageHandler{
		cfg,
	}
}

func (self MessagePageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(path.Join(self.cfg.HTMLDir, "messages.html"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("Could not load messages."))
		return
	}
	t.Execute(res, nil)
}
