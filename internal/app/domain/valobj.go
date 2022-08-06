package domain

type TeamPoints struct {
	Team   Team
	Points int
}

type TeamRatingChange struct {
	Team         Team
	RatingChange int
}
