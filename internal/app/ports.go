package app

import "ratingtable/internal/app/domain"

type UserFilter struct {
	Links map[string]string
}

type Storage interface {
	domain.RatingRepo

	// Game
	GameCreate(name string, gameType domain.GameType) (*domain.Game, error)
	GameGet(name string) (domain.IGame, error)

	// User
	UserCreate(links map[string]string) (*domain.User, error)
	UserFind(filter UserFilter) ([]domain.User, error)
	UserUpdateLinks(user domain.User, links map[string]string) error

	// Team
	TeamCreate(name string, gameID domain.GameID, users []domain.User) (*domain.Team, error)
	GetTeamParties(team domain.Team) ([]domain.Party, error)

	// Party
	PartyCreate(gameID domain.GameID, teamPoints []domain.TeamPoints) (*domain.Party, error)

	// Rating
	GetTableRating() ([]TeamRang, error)
}
