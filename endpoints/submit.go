package endpoints

import (
	"database/sql"
	"fmt"
	"net/http"

	"../config"
	"../teams"
)

// Submit handles POST requests to submit new flags and adjust team scores.
// Expects the following fields:
// 1. token - The submission token assigned to your team
// 2. flag  - The actual flag you are submitting
func Submit(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got a request to submit a flag")
		w.Header().Set("Content-Type", "text/plain")
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Your input is poorly formatted!"))
			return
		}
		fmt.Println(r.Form)
		tokens, found := r.Form["token"]
		if !found || len(tokens) == 0 {
			fmt.Println("Missing token")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing the token field. Please supply the submission token you were assigned."))
			return
		}
		flags, found := r.Form["flag"]
		if !found || len(flags) == 0 {
			fmt.Println("Missing flag")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing the flag field. Please supply secret flag."))
			return
		}
		team, err := teams.FindTeamByToken(db, tokens[0])
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("You submitted an invalid token. Please make sure you entered it correctly."))
			return
		}
		for _, flag := range cfg.Flags {
			if flag.Secret == flags[0] {
				team.Score += flag.Reward
				err = team.Update(db)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Could not update your score. Please notify the CTF administrators."))
					break
				} else {
					w.Write([]byte(fmt.Sprintf(
						"Congrats! You have been awarded %d points. Your score is now %d.\n",
						flag.Reward,
						team.Score)))
					break
				}
			}
		}
	}
}
