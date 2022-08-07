package app

import "ratingtable/internal/app/domain"

type Storage interface {
	domain.RatingRepo
	CreateGame(name string, _type domain.GameType) (*domain.Game, error)
}
