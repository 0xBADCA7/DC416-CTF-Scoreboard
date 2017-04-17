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

	QInitSessionsTable = `create table if not exists sessions (
		token varchar(16) primary key,
		created_at timestamp,
		expires_at timestamp
	);`

	QInitMessagesTable = `create table if not exists messages (
		id integer primary key,
		content text not null,
		created_at timestamp
	);`

	QSaveMessage = `
insert into messages (
	content, created_at
) values (
	?, ?
);`

	QAllMessages = `
select id, content, created_at
from messages
order by created_at desc;`

	QDeleteAllMessages = `delete from messages;`

	QGetTeams = `
select id, name, members, score, token, last_valid_submission
from teams;`

	QCreateTeam = `
insert into teams (
	name, members, score, token, last_valid_submission
) values (
	?, ?, 0, ?, datetime(0, 'unixepoch', 'localtime')
);`

	QCreateSession = `
insert into sessions (
	token, created_at, expires_at
) values (
	?, ?, ?
);`

	QDeleteSession = `
delete from sessions
where token = ?;`

	QFindTeamBySubmissionToken = `
select id, name, members, score, last_valid_submission
from teams
where token = ?;`

	QFindTeam = `
select name, token, members, score, last_valid_submission
from teams
where id = ?;`

	QDeleteTeam = `
delete from submitted where team_id = ?;
delete from teams where id = ?;`

	QUpdateTeam = `
update teams
set score = ?, token = ?, last_valid_submission = ?
where id = ?;`

	QFindSubmission = `
select id
from submitted
where team_id = ? and flag_id = ?;`

	QFindAllSubmissions = `
select id, flag_id
from submitted
where team_id = ?
order by id asc;`

	QFindSessionToken = `
select created_at, expires_at
from sessions
where token = ?;`

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
	_, err = db.Exec(QInitSessionsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(QInitMessagesTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(QInitSubmitted)
	return err
}

// generateUniqueToken creates a new 32-character hex-encoded string that is unique and can be used
// for things like authenticated session tokens and team submission tokens.
func generateUniqueToken(uniqueTest func(string) bool) string {
	buffer := make([]byte, 16)
	for {
		bytesRead, err := rand.Read(buffer)
		if err != nil || bytesRead != 16 {
			fmt.Println("Could not read random bytes for token.", err)
			continue
		}
		token := hex.EncodeToString(buffer)
		if uniqueTest(token) {
			return token
		}
	}
}
