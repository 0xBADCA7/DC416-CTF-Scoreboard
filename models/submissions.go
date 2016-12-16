package models

import (
	"database/sql"
)

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

// FindAllSubmissions attempts to find all of the submissions made by a given team.
func FindAllSubmissions(db *sql.DB, teamId int) ([]Submission, error) {
	submissions := []Submission{}
	rows, err := db.Query(QFindAllSubmissions, teamId)
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
