package app

import (
	"ratingtable/internal/app/domain"
)

type TeamRang struct {
	Rang   int
	Team   domain.Team
	Rating int
}

type App struct {
	storage Storage
}

func NewApp(storage Storage) *App {
	return &App{storage: storage}
}

func (a App) CreateGame(name string, gameType domain.GameType) error {
	_, err := a.storage.GameCreate(name, gameType)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetOrCreateUser(links map[string]string) (*domain.User, error) {
	users, err := a.storage.UserFind(UserFilter{Links: links})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return a.storage.UserCreate(links)
	}

	return &users[0], nil
}

func (a App) AddUserLinks(user domain.User, links map[string]string) error {
	return a.storage.UserUpdateLinks(user, links)
}

func (a App) AddTeam(name string, gameID domain.GameID, users []domain.User) (*domain.Team, error) {
	team, err := a.storage.TeamCreate(name, gameID, users)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (a App) AddParty(gameName string, teamsPoints []domain.TeamPoints) error {
	game, err := a.storage.GameGet(gameName)
	if err != nil {
		return err
	}

	err = game.AddParty(domain.Party{
		GameID:     game.GetID(),
		TeamPoints: teamsPoints,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetTableRating(users []domain.User) ([]TeamRang, error) {
	return a.storage.GetTableRating()
}

func (a App) GetTeamParties(team domain.Team) ([]domain.Party, error) {
	return a.storage.GetTeamParties(team)
}
