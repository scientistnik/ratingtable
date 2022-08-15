package domain

type TeamPoints struct {
	Team   Team
	Points float64
}

type TeamRatingChange struct {
	TeamID       TeamID
	PartyID      PartyID
	RatingChange int
}
