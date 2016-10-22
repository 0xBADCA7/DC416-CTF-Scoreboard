package main

import (
	"fmt"
	"net/http"

	"./config"
	"./endpoints"
)

func main() {
	cfg := config.Default()

	http.HandleFunc("/", endpoints.Index(nil, &cfg))
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
