package domain

type GameType int

const (
	_         GameType = iota
	ChessGame GameType = iota
)

type Game struct {
	ratingRepo RatingRepo
}

func (g Game) CalcRatingChanges(party Party) []TeamRatingChange {
	var teamRatings []TeamRatingChange

	var ratings []int
	for _, teamPoints := range party.TeamPoints {
		r := g.ratingRepo.getTeamRating(teamPoints.Team)
		ratings = append(ratings, r)
	}

	one, two := calcEloRating(
		ratings[0],
		ratings[1],
		party.TeamPoints[0].Points,
		party.TeamPoints[1].Points,
		40,
	)

	teamRatings = append(teamRatings, TeamRatingChange{Team: party.TeamPoints[0].Team, RatingChange: one})
	teamRatings = append(teamRatings, TeamRatingChange{Team: party.TeamPoints[1].Team, RatingChange: two})

	return teamRatings
}

func (g Game) AddParty(party Party) error {
	err := g.ratingRepo.addParty(party)
	if err != nil {
		return err
	}

	teamRatingChanges := g.CalcRatingChanges(party)

	err = g.ratingRepo.saveTeamRatingChanges(teamRatingChanges)
	if err != nil {
		return err
	}

	// TODO: update team ratings

	return nil
}
