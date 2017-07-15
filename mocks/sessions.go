package mocks

import (
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// SessionFindFn is the type of a function that effectively implements SessionModel.Find
type SessionFindFn func(string) (models.Session, error)

// SessionSaveFn is the type of a function that effectively implements SessionModel.Save
type SessionSaveFn func(*models.Session) error

// SessionSaveFn is the type of a function that effectively implements SessionModel.Delete
type SessionDeleteFn func(*models.Session) error

// SessionModelMock implements models.SessionMock in a way that lets us supply custom implementations of each method.
type SessionModelMock struct {
	find   SessionFindFn
	save   SessionSaveFn
	delete SessionDeleteFn
}

// NewSessionModelMock constructs a new mock implementation of models.SessionModel with caller-defined functions.
func NewSessionModelMock(find SessionFindFn, save SessionSaveFn, delete SessionDeleteFn) SessionModelMock {
	return SessionModelMock{
		find,
		save,
		delete,
	}
}

func (self SessionModelMock) Find(token string) (models.Session, error) {
	return self.find(token)
}

func (self SessionModelMock) Save(session *models.Session) error {
	return self.save(session)
}

func (self SessionModelMock) Delete(session *models.Session) error {
	return self.delete(session)
}
