package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Team contains information about teams that should never be served to users.
type Team struct {
	Id             int
	Name           string
	Members        string
	Score          int
	SubmitToken    string
	LastSubmission time.Time
}

// TeamScore contains information about a team's score.
// It is safe to serve a TeamScore to users.
type TeamScore struct {
	Name      string `json:"teamName"`  // The name of the team
	Members   string `json:"members"`   // The names of members, a single comma-separated string
	Score     int    `json:"score"`     // The team's score, as an integer
	Submitted string `json:"submitted"` // The time that the team's last submission occurred at
}

// TeamModel is implemented by types that can create and manage teams.
type TeamModel interface {
	Find(string) (Team, error)
	All() ([]Team, error)
	Save(*Team) error
	Update(*Team) error
	Delete(*Team) error
}

// TeamModelDB implements TeamModel such that Teams are stored in a sqlite database.
type TeamModelDB struct {
	db *sql.DB
}

// NewTeamModelDB constructs a new TeamModelDB with a database connection.
func NewTeamModelDB(db *sql.DB) TeamModelDB {
	return TeamModelDB{db}
}

// All obtains the current status of all of the teams.
func (self TeamModelDB) All() ([]Team, error) {
	rows, err := self.db.Query(QGetTeams)
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

// Find attempts to do a lookup for a team using its unique submission token.
func (self TeamModelDB) Find(token string) (Team, error) {
	team := Team{}
	fmt.Println("looking for team with token", token)
	err := self.db.QueryRow(QFindTeamBySubmissionToken, token).Scan(
		&team.Id, &team.Name, &team.Members, &team.Score, &team.LastSubmission)
	if err != nil {
		return Team{}, err
	}
	team.SubmitToken = token
	return team, err
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
func (self TeamModelDB) Save(team *Team) error {
	uniqueToken := generateUniqueToken(func(token string) bool {
		_, err := self.Find(token)
		return err != nil
	})
	team.SubmitToken = uniqueToken
	_, err := self.db.Exec(QCreateTeam, team.Name, team.Members, team.SubmitToken)
	if err != nil {
		return err
	}
	err = self.db.QueryRow(QLastInsertedId).Scan(team.Id)
	return err
}

// Update resets the team's score and allows for changing their submission token.
func (self TeamModelDB) Update(team *Team) error {
	_, err := self.db.Exec(QUpdateTeam, team.Score, team.SubmitToken, team.LastSubmission, team.Id)
	return err
}

// Delete removes a team from the database.
func (self TeamModelDB) Delete(team *Team) error {
	_, err := self.db.Exec(QDeleteTeam, team.Id, team.Id)
	team.Id = -1
	return err
}
