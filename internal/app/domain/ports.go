package domain

type RatingRepo interface {
	addParty(party Party) error
	saveTeamRatingChanges([]TeamRatingChange) error
	getTeamRating(team Team) int
}
