package domain

type Chess struct {
	Game
}

func NewChessGame(repo RatingRepo) *Chess {
	return &Chess{Game{repo}}
}
