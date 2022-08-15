package domain

type UserID int

type User struct {
	ID UserID
}

type TeamID int

type Team struct {
	ID    TeamID
	Users []User
}

type GameID int

type Game struct {
	ID   GameID
	Name string
	Type GameType
}

type PartyID int

type Party struct {
	ID         PartyID
	GameID     GameID
	TeamPoints []TeamPoints
}
