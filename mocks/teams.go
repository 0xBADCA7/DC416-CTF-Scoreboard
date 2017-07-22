package mocks

import (
	"github.com/DC416/DC416-CTF-Scoreboard/models"
)

// TeamFindFn is the type of a function that effectively satisfies TeamModel.Find
type TeamFindFn func(token string) (models.Team, error)

// TeamAllFn is the type of a function that effectively satisfies TeamModel.All
type TeamAllFn func() ([]models.Team, error)

// TeamSaveFn is the type of a function that effectively satisfies TeamModel.Save
type TeamSaveFn func(*models.Team) error

// TeamUpdateFn is the type of a function that effectively satisfies TeamModel.Update
type TeamUpdateFn func(*models.Team) error

// TeamDeleteFn is the type of a function that effectively satisfies TeamModel.Delete
type TeamDeleteFn func(*models.Team) error

// TeamModelMock implements models.TeamModel in a way that lets us supply custom implementations of each method.
type TeamModelMock struct {
	find   TeamFindFn
	all    TeamAllFn
	save   TeamSaveFn
	update TeamUpdateFn
	delete TeamDeleteFn
}

// NewTeamModelMock constructs a new mock implementation of models.TeamModel with caller-defined functions.
func NewTeamModelMock(
	find TeamFindFn,
	all TeamAllFn,
	save TeamSaveFn,
	update TeamUpdateFn,
	delete TeamDeleteFn) TeamModelMock {
	return TeamModelMock{
		find,
		all,
		save,
		update,
		delete,
	}
}

func (self TeamModelMock) Find(token string) (models.Team, error) {
	return self.find(token)
}

func (self TeamModelMock) All() ([]models.Team, error) {
	return self.all()
}

func (self TeamModelMock) Save(team *models.Team) error {
	return self.save(team)
}

func (self TeamModelMock) Update(team *models.Team) error {
	return self.update(team)
}

func (self TeamModelMock) Delete(team *models.Team) error {
	return self.delete(team)
}
