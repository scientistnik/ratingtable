package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"ratingtable/internal/app"
	"ratingtable/internal/app/domain"

	packr "github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
)

var migrations *migrate.PackrMigrationSource = &migrate.PackrMigrationSource{
	Box: packr.New("migrations", "./migrations"),
}

type SQLite struct {
	filename string
	db       *sql.DB
}

var _ app.Storage = (*SQLite)(nil)

func NewSQLiteStorage(filename string) (*SQLite, error) {
	instance := SQLite{filename: filename}
	err := instance.migrateRollUp()
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (s SQLite) migrateRollUp() error {
	db, err := sql.Open("sqlite3", s.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return err
	}

	if n > 0 {
		fmt.Printf("Applied %d migrations!\n", n)
	}

	return nil
}

func (s *SQLite) Connect() error {
	db, err := sql.Open("sqlite3", s.filename)
	if err != nil {
		return err
	}

	s.db = db

	return err
}

func (s SQLite) Disconnect() error {
	return s.db.Close()
}

func (s SQLite) AddParty(party domain.Party) error {

	result, err := s.db.Exec("INSERT INTO parties (game_id, created_at) VALUES (?,?)", party.GameID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, teamPoints := range party.TeamPoints {
		_, err = s.db.Exec("INSERT INTO party_teampoints (party_id, team_id, points) VALUES (?,?,?)", id, teamPoints.Team.ID, teamPoints.Points)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SQLite) SaveTeamRatingChanges(ratingChanges []domain.TeamRatingChange) error {
	// TODO: check what updated_at is last
	for _, ratingChange := range ratingChanges {
		_, err := s.db.Exec(
			"INSERT INTO rating_changes (team_id, party_id, points_change) VALUES (?,?,?)",
			ratingChange.TeamID,
			ratingChange.PartyID,
			ratingChange.RatingChange,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SQLite) GetTeamRating(team domain.Team) (int, error) {
	rows, err := s.db.Query("SELECT id, rating from ratings where team_id = ?", team.ID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, errors.New("not found team")
	}

	var rating int
	err = rows.Scan(&rating)
	if err != nil {
		return 0, err
	}

	return rating, nil
}

func (s SQLite) RecalcTeamRating(team domain.Team) error {
	return nil
}

func (s SQLite) GameCreate(name string, gameType domain.GameType) (*domain.Game, error) {
	result, err := s.db.Exec("INSERT INTO games (name, type) VALUES (?,?)", name, gameType)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &domain.Game{ID: domain.GameID(id), Name: name, Type: gameType}, nil
}

func (s SQLite) GameGet(name string) (domain.IGame, error) {
	rows, err := s.db.Query("SELECT id, type from games where name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("not found game")
	}

	var id domain.GameID
	var gameType domain.GameType
	err = rows.Scan(&id, &name, &gameType)
	if err != nil {
		return nil, err
	}

	return domain.DefaultGame{ID: id, Name: name}, nil
}

func (s SQLite) UserCreate(links map[string]string) (*domain.User, error) {
	return nil, nil
}

func (s SQLite) UserFind(filter app.UserFilter) ([]domain.User, error) {
	return nil, nil
}

func (s SQLite) UserUpdateLinks(user domain.User, links map[string]string) error {
	return nil
}

func (s SQLite) TeamCreate(name string, gameID domain.GameID, users []domain.User) (*domain.Team, error) {
	return nil, nil
}

func (s SQLite) GetTeamParties(team domain.Team) ([]domain.Party, error) {
	return nil, nil
}

func (s SQLite) PartyCreate(gameID domain.GameID, teamPoints []domain.TeamPoints) (*domain.Party, error) {
	return nil, nil
}

func (s SQLite) GetTableRating() []app.TeamRang {
	return nil
}
