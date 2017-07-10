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

	sessions := models.NewSessionModel(db)
	submissions := models.NewSubmissionModelDB(db)
	teams := models.NewTeamModelDB(db)

	submissionHandler := endpoints.NewSubmissionHandler(cfg, submissions, teams)
	registrationHandler := endpoints.NewRegistrationHandler(cfg, teams, sessions)
	adminPageHandler := endpoints.NewAdminPageHandler(cfg, submissions, teams, sessions)
	deleteTeamHandler := endpoints.NewDeleteTeamHandler(teams, sessions)

	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/img/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", endpoints.Index(db, &cfg))
	http.Handle("/register", registrationHandler)
	http.Handle("/submit", submissionHandler)
	http.HandleFunc("/login", endpoints.Login(db, &cfg))
	http.HandleFunc("/logout", endpoints.Logout(db, &cfg))
	http.Handle("/admin", adminPageHandler)
	http.HandleFunc("/message", endpoints.PostMessage(db, &cfg))
	http.Handle("/deleteteam", deleteTeamHandler)
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
