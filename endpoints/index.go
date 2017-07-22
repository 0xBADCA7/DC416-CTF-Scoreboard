package endpoints

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type IndexHandler struct {
	cfg   config.Config
	teams models.TeamModel
}

func NewIndexHandler(cfg config.Config, teams models.TeamModel) IndexHandler {
	return IndexHandler{
		cfg,
		teams,
	}
}

// Index creates a request handler that serves index.html, the main scoreboard page with all of
// the teams and their scores.
func (self IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templateFile, err := os.Open(path.Join(self.cfg.HTMLDir, "index.html"))
	if err != nil {
		fmt.Println("Error parsing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("We done goofed! Try again in a few minutes."))
		return
	}
	templateCode, readErr := ioutil.ReadAll(templateFile)
	if readErr != nil {
		fmt.Println("Error reading template for index.html", readErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("We done goofed! Try again in a few minutes."))
		return
	}
	t, err := template.New("index").Parse(string(templateCode))
	if err != nil {
		fmt.Println("Error building template", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("We done goofed! Try again in a few minutes."))
		return
	}
	t.Execute(w, nil)
}
