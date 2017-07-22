package mocks

import (
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// SubmissionFindFn is the type of a function that effectively satisfies SubmissionModel.Save
type SubmissionFindFn func(int, int) (models.Submission, error)

// SubmissionAllFn is the type of a function that effectively satisfies SubmissionModel.All
type SubmissionAllFn func() ([]models.Submission, error)

// SubmissionSaveFn is the type of a function that effectively satisfies SubmissionModel.Save
type SubmissionSaveFn func(*models.Submission) error

// SubmissionModelMock implements models.SubmissionModel in a way that lets us supply custom implementations of each method.
type SubmissionModelMock struct {
	find SubmissionFindFn
	all  SubmissionAllFn
	save SubmissionSaveFn
}

// NewSubmissionModelMock constructs a new mock implementation of models.SubmissionModel with caller-defined functions.
func NewSubmissionModelMock(find SubmissionFindFn, all SubmissionAllFn, save SubmissionSaveFn) SubmissionModelMock {
	return SubmissionModelMock{
		find,
		all,
		save,
	}
}

func (self SubmissionModelMock) Find(teamId, flagId int) (models.Submission, error) {
	return self.find(teamId, flagId)
}

func (self SubmissionModelMock) All() ([]models.Submission, error) {
	return self.all()
}

func (self SubmissionModelMock) Save(submission *models.Submission) error {
	return self.save(submission)
}
