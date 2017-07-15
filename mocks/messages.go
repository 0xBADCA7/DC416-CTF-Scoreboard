package mocks

import (
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// MessageSaveFn is the type of a function that effectively satisfies MessageModel.Save
type MessageSaveFn func(*models.Message) error

// MessageAllFn is the type of a function that effectively satisfies MessageModel.All
type MessageAllFn func() ([]models.Message, error)

// MessageDeleteFn is the type of a function that effectively satisfies MessageModel.Delete
type MessageDeleteFn func(*models.Message) error

// MessageModelMock implements models.MessageModel in a way that lets us supply custom implementations of each method.
type MessageModelMock struct {
	save   MessageSaveFn
	all    MessageAllFn
	delete MessageDeleteFn
}

// NewMessageModelMock constructs a new mock implementation of models.MessageModel with caller-defined functions.
func NewMessageModelMock(save MessageSaveFn, all MessageAllFn, delete MessageDeleteFn) MessageModelMock {
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
