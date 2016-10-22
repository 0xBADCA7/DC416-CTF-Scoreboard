package config

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
	HTMLDir      string `json:"htmlDir"`     // The path to the directory housing HTML files, relative to main.go
	CSSDir       string `json:"cssDir"`      // The path to the directory housing CSS files, relative to main.go
	JSDir        string `json:"jsDir"`       // The path to the directory housing JavaScript files, relative to main.go
	ImgDir       string `json:"imgDir"`      // The path to the directory housing image files, relative to main.go
	Flags        []Flag `json:"flags"`       // Information about flags that users can submit to get points
}

func Default() Config {
	return Config{
		"localhost:3000",
		"scoreboard.db",
		"html",
		"css",
		"js",
		"img",
		[]Flag{
			{1, "secret1", 10},
			{2, "secret2", 20},
		},
	}
}
