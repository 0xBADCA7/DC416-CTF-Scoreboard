package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// FindTeams obtains the current status of all of the teams.
func FindTeams(db *sql.DB) ([]Team, error) {
	rows, err := db.Query(QGetTeams)
	if err != nil {
		return []Team{}, err
	}
	teams := []Team{}
	for rows.Next() {
		team := Team{}
		err = rows.Scan(&team.Id, &team.Name, &team.Members, &team.Score, &team.SubmitToken, &team.LastSubmission)
		if err == nil {
			fmt.Println("Got team", team)
			teams = append(teams, team)
		} else {
			fmt.Println(err)
		}
	}
	return teams, err
}

// FindTeamByToken attempts to do a lookup for a team using its unique submission token.
func FindTeamByToken(db *sql.DB, token string) (Team, error) {
	team := Team{}
	fmt.Println("looking for team with token", token)
	err := db.QueryRow(QFindTeamBySubmissionToken, token).Scan(
		&team.Id, &team.Name, &team.Members, &team.Score, &team.LastSubmission)
	if err != nil {
		return Team{}, err
	}
	team.SubmitToken = token
	return team, err
}

// generateUniqueToken creates a new 32-character hex-encoded string that is unique and can
// be used as a security token by teams submitting flags.
func generateUniqueToken(db *sql.DB) string {
	buffer := make([]byte, 16)
	for {
		bytesRead, err := rand.Read(buffer)
		if err != nil || bytesRead != 16 {
			fmt.Println("Could not read random bytes for token.", err)
			continue
		}
		token := hex.EncodeToString(buffer)
		_, err = FindTeamByToken(db, token)
		if err != nil {
			return token
		}
	}
}

// Team contains information about teams that should never be served to users.
type Team struct {
	Id             int
	Name           string
	Members        string
	Score          int
	SubmitToken    string
	LastSubmission time.Time
}

// TeamByScore is an alias for an array of teams that implements all of the interfaces required by
// the sort package to be able to sort teams by their score.
type TeamByScore []Team

func (t TeamByScore) Len() int {
	return len(t)
}

func (t TeamByScore) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TeamByScore) Less(i, j int) bool {
	if t[i].Score == t[j].Score {
		return t[j].LastSubmission.After(t[i].LastSubmission)
	}
	return t[i].Score > t[j].Score
}

// Save creates a new Team in the database.
func (t *Team) Save(db *sql.DB) error {
	uniqueToken := generateUniqueToken(db)
	t.SubmitToken = uniqueToken
	_, err := db.Exec(QCreateTeam, t.Name, t.Members, t.SubmitToken)
	fmt.Println("---", err)
	return err
}

// Update resets the team's score and allows for changing their submission token.
func (t *Team) Update(db *sql.DB) error {
	_, err := db.Exec(QUpdateTeam, t.Score, t.SubmitToken, t.LastSubmission, t.Id)
	return err
}

// TeamScore contains information about a team's score.
// It is safe to serve a TeamScore to users.
type TeamScore struct {
	Name      string `json:"teamName"`  // The name of the team
	Members   string `json:"members"`   // The names of members, a single comma-separated string
	Score     int    `json:"score"`     // The team's score, as an integer
	Submitted string `json:"submitted"` // The time that the team's last submission occurred at
}
