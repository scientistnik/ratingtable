package domain

type RatingRepo interface {
	AddParty(party Party) error
	SaveTeamRatingChanges([]TeamRatingChange) error
	GetTeamRating(team Team) (int, error)
	RecalcTeamRating(team Team) error
}
