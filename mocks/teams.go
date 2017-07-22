package mocks

import (
	"errors"

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

// Closed over by the TeamModelMock created by NewInMemoryTeamModel.
type teamState struct {
	Teams []models.Team
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

// Constructs a TeamModelMock that operates on an in-memory array of teams
// instead of talking to a database.
//
// Note that this mock has the quirk that, after a team with an id N gets
// deleted, all teams with ids m >= n will be decremented to id = m - 1.
func NewInMemoryTeamModel() TeamModelMock {
	state := teamState{make([]models.Team, 0)}

	find := func(token string) (models.Team, error) {
		for _, team := range state.Teams {
			if team.SubmitToken == token {
				return team, nil
			}
		}
		team := models.Team{}
		return team, errors.New("Team not found")
	}

	all := func() ([]models.Team, error) {
		return state.Teams, nil
	}

	save := func(team *models.Team) error {
		team.Id = len(state.Teams) + 1
		state.Teams = append(state.Teams, *team)
		return nil
	}

	update := func(team *models.Team) error {
		for i := 0; i < len(state.Teams); i++ {
			if state.Teams[i].Id == team.Id {
				state.Teams[i] = *team
				return nil
			}
		}
		return errors.New("Team not found")
	}

	del := func(team *models.Team) error {
		index := -1
		for i := 0; i < len(state.Teams); i++ {
			if state.Teams[i].Id == team.Id {
				index = i
				break
			}
		}
		if index < 0 {
			return errors.New("Team not found")
		}
		state.Teams = append(state.Teams[:index], state.Teams[index+1:]...)
		return nil
	}

	return NewTeamModelMock(find, all, save, update, del)
}
