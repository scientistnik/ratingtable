package domain

type User struct {
	ID int
}

type Team struct {
	ID int
}

type Party struct {
	ID         int
	GameID     int
	TeamPoints []TeamPoints
}
