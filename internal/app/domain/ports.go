package domain

type TeamPoints struct {
	Team   Team
	Points int
}

type Game interface {
	AddPartyResult([]TeamPoints) error
}

type RatingRepo interface {
	addParty(gameID int) error
	addPartyTeam(partyID int, teamID int, points int) error
}
