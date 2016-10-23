package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"./config"
	"./endpoints"
	"./teams"
)

func main() {
	cfg := config.MustLoad("./config/config.json")

	db, err := sql.Open("sqlite3", cfg.DatabaseFile)
	if err != nil {
		panic(err)
	}
	err = teams.InitTables(db)
	if err != nil {
		panic(err)
	}

	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/img/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", endpoints.Index(db, &cfg))
	http.HandleFunc("/register", endpoints.Register(db, &cfg))
	http.HandleFunc("/submit", endpoints.Submit(db, &cfg))
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
