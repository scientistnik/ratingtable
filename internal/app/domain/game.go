package domain

type GameType int

const (
	DefaultGameType GameType = iota
	ChessGameType   GameType = iota
)

type IGame interface {
	GetID() GameID
	AddParty(party Party) error
}

type DefaultGame struct {
	ID         GameID
	Name       string
	ratingRepo RatingRepo
}

var _ IGame = (*DefaultGame)(nil)

func (g DefaultGame) GetID() GameID {
	return g.ID
}

func (g DefaultGame) CalcRatingChanges(party Party) []TeamRatingChange {
	var teamRatings []TeamRatingChange

	var ratings []int
	for _, teamPoints := range party.TeamPoints {
		r, err := g.ratingRepo.GetTeamRating(teamPoints.Team)
		if err != nil {
			return nil
		}

		ratings = append(ratings, r)
	}

	one := calcEloRating(
		ratings[0],
		float64(party.TeamPoints[0].Points),
		ratings[1],
		float64(party.TeamPoints[1].Points),
		40,
	)
	two := calcEloRating(
		ratings[1],
		float64(party.TeamPoints[1].Points),
		ratings[0],
		float64(party.TeamPoints[0].Points),
		40,
	)

	teamRatings = append(teamRatings, TeamRatingChange{TeamID: party.TeamPoints[0].Team.ID, RatingChange: one, PartyID: party.ID})
	teamRatings = append(teamRatings, TeamRatingChange{TeamID: party.TeamPoints[1].Team.ID, RatingChange: two, PartyID: party.ID})

	return teamRatings
}

func (g DefaultGame) AddParty(party Party) error {
	err := g.ratingRepo.AddParty(party)
	if err != nil {
		return err
	}

	teamRatingChanges := g.CalcRatingChanges(party)

	err = g.ratingRepo.SaveTeamRatingChanges(teamRatingChanges)
	if err != nil {
		return err
	}

	for _, teamPoints := range party.TeamPoints {
		err := g.ratingRepo.RecalcTeamRating(teamPoints.Team)
		if err != nil {
			return err
		}
	}

	return nil
}
