package config

import (
	"encoding/json"
	"os"
)

// PasswordEnvVar is the name of the environment variable that we expect the password for admin access
// to be provided under.
const PasswordEnvVar = "IVS_PASSWORD"

// Flag contains a unique identifier
type Flag struct {
	Id     int    `json:"id"`
	Secret string `json:"secret"`
	Reward int    `json:"reward"`
}

// Config contains all of the configuration information required by the application.
type Config struct {
	BindAddress  string `json:"bindAddress"` // The address to bind the server to, formatted <ip>:<port>
	DatabaseFile string `json:"dbFile"`      // The name of the file to have SQLite save the database state to
	CTFName      string `json:"ctfName"`     // The name of the CTF event, to display over the scoreboard
	HTMLDir      string `json:"htmlDir"`     // The path to the directory housing HTML files, relative to main.go
	Flags        []Flag `json:"flags"`       // Information about flags that users can submit to get points
}

// MustLoad will try to load a JSON file containing a valid configuration for the application.
// Failing to do so, it will panic, terminating the program.
func MustLoad(configFilePath string) Config {
	cfg := Config{}
	f, openErr := os.Open(configFilePath)
	if openErr != nil {
		panic(openErr)
	}
	decoder := json.NewDecoder(f)
	decodeErr := decoder.Decode(&cfg)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return cfg
}
