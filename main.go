package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"./authentication"
	"./config"
	"./endpoints"
	"./models"
)

func main() {
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

	authentication.HashAdminPassword()

	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/img/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", endpoints.Index(db, &cfg))
	http.HandleFunc("/register", endpoints.Register(db, &cfg))
	http.HandleFunc("/submit", endpoints.Submit(db, &cfg))
	http.HandleFunc("/login", endpoints.Login(db, &cfg))
	http.HandleFunc("/logout", endpoints.Logout(db, &cfg))
	http.HandleFunc("/admin", endpoints.Admin(db, &cfg))
	http.HandleFunc("/message", endpoints.PostMessage(db, &cfg))
	http.HandleFunc("/deleteteam", endpoints.DeleteTeam(db, &cfg))
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
