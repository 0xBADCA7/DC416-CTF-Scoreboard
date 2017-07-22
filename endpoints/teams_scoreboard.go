package endpoints

import (
	"net/http"

	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

type TeamScoreboardHandler struct {
	teams models.TeamModel
}

func NewTeamsScoreboardHandler(teams models.TeamModel) TeamScoreboardHandler {
	return TeamScoreboardHandler{
		teams,
	}
}

func (self TeamScoreboardHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("fail"))
}
