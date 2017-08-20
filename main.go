package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"github.com/DC416/DC416-CTF-Scoreboard/authentication"
	"github.com/DC416/DC416-CTF-Scoreboard/config"
	"github.com/DC416/DC416-CTF-Scoreboard/endpoints"
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

func main() {
	authentication.HashAdminPassword()

	cfgFile := os.Getenv("CONFIG_FILE")
	if cfgFile == "" {
		cfgFile = "./config/config.json"
	}
	fmt.Println(cfgFile)
	cfg := config.MustLoad(cfgFile)

	db, err := sql.Open("sqlite3", cfg.DatabaseFile)
	if err != nil {
		panic(err)
	}
	err = models.InitTables(db)
	if err != nil {
		panic(err)
	}

	messages := models.NewMessageModelDB(db)
	teams := models.NewTeamModelDB(db)
	sessions := models.NewSessionModelDB(db)
	submissions := models.NewSubmissionModelDB(db)

	registrationHandler := endpoints.NewRegistrationHandler(cfg, teams, sessions)
	deleteTeamHandler := endpoints.NewDeleteTeamHandler(teams, sessions)

	indexHandler := endpoints.NewIndexHandler(cfg, teams)
	eventInfoHandler := endpoints.NewEventInfoHandler(cfg.CTFName)
	scoreboardHandler := endpoints.NewTeamsScoreboardHandler(teams)
	submitPageHandler := endpoints.NewSubmitPageHandler(cfg)
	submissionHandler := endpoints.NewTeamSubmitHandler(teams, submissions, cfg.Flags)
	loginPageHandler := endpoints.NewLoginPageHandler(cfg)
	loginHandler := endpoints.NewLoginHandler(sessions)
	messagesHandler := endpoints.NewMessagesHandler(messages)
	messagePageHandler := endpoints.NewMessagePageHandler(cfg)
	messagePostHandler := endpoints.NewMessagesPostHandler(messages, sessions)
	adminPageHandler := endpoints.NewAdminPageHandler(cfg, submissions, teams, sessions)
	adminTeamsHandler := endpoints.NewAdminTeamsHandler(cfg, submissions, teams, sessions)

	router := mux.NewRouter()

	http.Handle("/css/", http.FileServer(http.Dir(".")))
	// http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/img/", http.FileServer(http.Dir(".")))
	router.Handle("/", indexHandler)
	router.Handle("/event", eventInfoHandler)
	router.Handle("/teams/scoreboard", scoreboardHandler)
	router.Handle("/register", registrationHandler)
	router.Handle("/submit", submitPageHandler).Methods("GET")
	router.Handle("/submit", submissionHandler).Methods("POST")
	router.Handle("/login", loginPageHandler).Methods("GET")
	router.Handle("/login", loginHandler).Methods("POST")
	router.Handle("/admin", adminPageHandler).Methods("GET")
	router.Handle("/admin/teams", adminTeamsHandler).Methods("GET")
	router.Handle("/admin/messages", messagePostHandler).Methods("POST")
	router.Handle("/message", messagePageHandler).Methods("GET")
	router.Handle("/messages", messagesHandler).Methods("GET")
	router.Handle("/deleteteam", deleteTeamHandler)

	http.Handle("/", router)
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
