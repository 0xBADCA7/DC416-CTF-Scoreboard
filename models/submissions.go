package models

import (
	"database/sql"
)

// Submission contains information about a flag submitted by a team.
// The flag ID is an integer that is specified in the application config and is unique.
type Submission struct {
	Id    int
	Flag  int
	Owner int
}

// SubmissionModel is implemented by types that can lookup submissions made by teams
// and store new submissions for teams.
type SubmissionModel interface {
	Find(int, int) (Submission, error)
	All(int) ([]Submission, error)
	Save(*Submission) error
}

// SubmissionModelDB implements SubmissionModel to work with a sqlite datbase.
type SubmissionModelDB struct {
	db *sql.DB
}

// NewSubmission creates a new Submission for a given flag being sent by a given team.
func NewSubmission(flagId, ownerId int) Submission {
	return Submission{
		-1,
		flagId,
		ownerId,
	}
}

// NewSubmissionModelDB constructs a new SubmissionModelDB capable of working with a sqlite database.
func NewSubmissionModelDB(db *sql.DB) SubmissionModelDB {
	return SubmissionModelDB{db}
}

// Find attempts to find an entry in the submission table for the flag that the user is submitting
// for their team.
func (self SubmissionModelDB) Find(teamId int, flagId int) (Submission, error) {
	s := Submission{}
	err := self.db.QueryRow(QFindSubmission, teamId, flagId).Scan(&s.Id)
	if err != nil {
		return Submission{}, err
	}
	s.Flag = flagId
	s.Owner = teamId
	return s, nil
}

// All attempts to find all of the submissions made by a given team.
func (self SubmissionModelDB) All(teamId int) ([]Submission, error) {
	submissions := []Submission{}
	rows, err := self.db.Query(QFindAllSubmissions, teamId)
	if err != nil {
		return submissions, err
	}
	for rows.Next() {
		var id, flagId int
		err = rows.Scan(&id, &flagId)
		if err != nil {
			return []Submission{}, err
		}
		submissions = append(submissions, Submission{
			Id:    id,
			Flag:  flagId,
			Owner: teamId,
		})
	}
	return submissions, nil
}

// Save creates a new record of a team submitting a flag.
func (self SubmissionModelDB) Save(submission *Submission) error {
	_, err := self.db.Exec(QSaveSubmission, submission.Owner, submission.Flag)
	if err != nil {
		return err
	}
	err = self.db.QueryRow(QLastInsertedId).Scan(submission.Id)
	return err
}
