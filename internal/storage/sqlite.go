package storage

import (
	"database/sql"
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

func (rr SQLite) AddParty(party domain.Party) error {
	return nil
}

func (rr SQLite) SaveTeamRatingChanges([]domain.TeamRatingChange) error {
	return nil
}

func (rr SQLite) GetTeamRating(team domain.Team) int {
	return 0
}

func (rr SQLite) RecalcTeamRating(team domain.Team) error {
	return nil
}

func (rr SQLite) CreateGame(name string, gameType domain.GameType) (*domain.Game, error) {
	return nil, nil
}
