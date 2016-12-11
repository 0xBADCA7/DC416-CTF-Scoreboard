package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
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
