package mocks

import (
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// SaveFn is the type of a function that effectively satisfies MessageModel.Save
type SaveFn func(*models.Message) error

// AllFn is the type of a function that effectively satisfies MessageModel.All
type AllFn func() ([]models.Message, error)

// DeleteFn is the type of a function that effectively satisfies MessageModel.Delete
type DeleteFn func(*models.Message) error

// MessageModelMock implements models.MessageModel in a way that lets us supply custom implementations of each method.
type MessageModelMock struct {
	save   SaveFn
	all    AllFn
	delete DeleteFn
}

// NewMessageModelMock constructs a new mock implementation of models.MessageModel with caller-defined functions.
func NewMessageModelMock(save SaveFn, all AllFn, delete DeleteFn) MessageModelMock {
	return MessageModelMock{
		save,
		all,
		delete,
	}
}

func (self MessageModelMock) Save(message *models.Message) error {
	return self.save(message)
}

func (self MessageModelMock) All() ([]models.Message, error) {
	return self.all()
}

func (self MessageModelMock) Delete(message *models.Message) error {
	return self.delete(message)
}
