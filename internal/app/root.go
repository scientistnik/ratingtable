package app

import "ratingtable/internal/app/domain"

type UserRang struct {
	Rang   int
	User   domain.User
	Rating int
}

type App struct {
	storage Storage
}

func (a App) CreateGame(name string, gameType domain.GameType) error {

	_, err := a.storage.CreateGame(name, gameType)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetOrCreateUser(links []string) (*domain.User, error) {
	return nil, nil
}

func (a App) AddUserLinks(links []string) error {
	return nil
}

func (a App) AddTeam(name string, gameID int, users []domain.User) (*domain.Team, error) {
	return nil, nil
}

func (a App) AddUserTeam(team domain.Team, users []domain.User) error {
	return nil
}

func (a App) AddParty(party domain.Party) error {
	return nil
}

func (a App) AddPartyByLink(link string) error {
	return nil
}

func (a App) AddPartyByResultStr(resultStr string) error {
	return nil
}

func (a App) GetUsersRating(users []domain.User) []int {
	return nil
}

func (a App) GetTableRating() []UserRang {
	return nil
}
