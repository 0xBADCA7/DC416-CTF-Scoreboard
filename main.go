package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

	submissionHandler := endpoints.NewSubmissionHandler(cfg, submissions, teams)
	registrationHandler := endpoints.NewRegistrationHandler(cfg, teams, sessions)
	adminPageHandler := endpoints.NewAdminPageHandler(cfg, submissions, teams, sessions)
	deleteTeamHandler := endpoints.NewDeleteTeamHandler(teams, sessions)
	messageHander := endpoints.NewMessageHandler(cfg, messages, sessions)
	loginHandler := endpoints.NewLoginHandler(cfg, sessions)

	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/img/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", endpoints.Index(db, &cfg))
	http.Handle("/register", registrationHandler)
	http.Handle("/submit", submissionHandler)
	http.Handle("/login", loginHandler)
	http.HandleFunc("/logout", endpoints.Logout(db, &cfg))
	http.Handle("/admin", adminPageHandler)
	http.Handle("/message", messageHander)
	http.Handle("/deleteteam", deleteTeamHandler)
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
