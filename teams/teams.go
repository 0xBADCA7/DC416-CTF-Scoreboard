package teams

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	QInitTeamTable = `create table if not exists teams (
		id integer primary key,
		name text unique,
		members text,
		score integer,
		token text unique,
		last_valid_submission timestamp
);`

	QInitSubmitted = `create table if not exists submitted (
		id integer primary key,
		team_id integer,
		flag_id integer,
		foreign key(team_id) references teams(id)
	);`

	QGetTeams = `
select id, name, members, score, token, last_valid_submission
from teams;`

	QCreateTeam = `
insert into teams (
	name, members, score, token, last_valid_submission
) values (
	?, ?, 0, ?, datetime(0, 'unixepoch', 'localtime')
);`

	QFindTeamBySubmissionToken = `
select id, name, members, score, last_valid_submission
from teams
where token = ?;`

	QUpdateTeam = `
update teams
set score = ?, token = ?, last_valid_submission = ?
where id = ?;`

	QFindSubmission = `
select id
from submitted
where team_id = ? and flag_id = ?;`

	QSaveSubmission = `
insert into submitted (
	team_id, flag_id
) values (
	?, ?
);`
)

// InitTables initializes the database tables.
func InitTables(db *sql.DB) error {
	_, err := db.Exec(QInitTeamTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(QInitSubmitted)
	return err
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

// FindSubmission attempts to find an entry in the submission table for the flag that the user is submitting
// for their team.
func FindSubmission(db *sql.DB, teamId int, flagId int) (Submission, error) {
	s := Submission{}
	err := db.QueryRow(QFindSubmission, teamId, flagId).Scan(&s.Id)
	if err != nil {
		return Submission{}, err
	}
	s.Flag = flagId
	s.Owner = teamId
	return s, nil
}

// Submission contains information about a flag submitted by a team.
// The flag ID is an integer that is specified in the application config and is unique.
type Submission struct {
	Id    int
	Flag  int
	Owner int
}

// Save creates a new record of a team submitting a flag.
func (s *Submission) Save(db *sql.DB) error {
	_, err := db.Exec(QSaveSubmission, s.Owner, s.Flag)
	return err
}

// TeamScore contains information about a team's score.
// It is safe to serve a TeamScore to users.
type TeamScore struct {
	Name    string `json:"teamName"` // The name of the team
	Members string `json:"members"`  // The names of members, a single comma-separated string
	Score   int    `json:"score"`    // The team's score, as an integer
}
