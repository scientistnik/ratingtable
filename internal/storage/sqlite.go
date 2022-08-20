package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ratingtable/internal/app"
	"ratingtable/internal/app/domain"
	"strings"

	_ "github.com/mattn/go-sqlite3"

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
	// TODO
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

	return domain.GetGameByType(gameType, s, id, name), nil
}

func (s SQLite) UserCreate(links map[string]string) (*domain.User, error) {
	linksJson, err := json.Marshal(links)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec("INSERT INTO users (links) VALUES (?)", linksJson)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: domain.UserID(id)}, nil
}

func (s SQLite) UserFind(filter app.UserFilter) ([]domain.User, error) {
	var predicats []string

	if filter.Links != nil {
		var linkPredicats []string

		for key, value := range filter.Links {
			linkPredicats = append(linkPredicats, fmt.Sprintf("json_extract(links,'%s') = '%s'", key, value))
		}

		predicats = append(predicats, "("+strings.Join(linkPredicats, " or ")+")")
	}

	rows, err := s.db.Query("SELECT id from users where " + strings.Join(predicats, " and "))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s SQLite) UserUpdateLinks(user domain.User, links map[string]string) error {
	rows, err := s.db.Query("SELECT links FROM users WHERE id = ?", user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("not found user(id=%d)", user.ID)
	}

	var currLinks map[string]string
	err = rows.Scan(&currLinks)
	if err != nil {
		return err
	}

	for key, value := range links {
		currLinks[key] = value
	}

	_, err = s.db.Exec("UPDATE users set links = ? WHERE id = ?", currLinks, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s SQLite) TeamCreate(name string, gameID domain.GameID, users []domain.User) (*domain.Team, error) {
	result, err := s.db.Exec("INSERT INTO teams (name, game_id) VALUES (?,?)", name, gameID)
	if err != nil {
		return nil, err
	}

	teamID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var values []any
	var query []string
	for _, user := range users {
		query = append(query, "(?,?)")
		values = append(values, user.ID, teamID)
	}

	_, err = s.db.Exec("INSERT INTO user_team (user_id, team_id) VALUES "+strings.Join(query, ","), values...)
	if err != nil {
		return nil, err
	}

	return &domain.Team{ID: domain.TeamID(teamID), Name: name, Users: users}, nil
}

func (s SQLite) GetTeamParties(team domain.Team) ([]domain.Party, error) {
	// TODO: select all teams in party
	rows, err := s.db.Query("SELECT p.id, p.game_id, p.created_at, ptp.points FROM parties as p "+
		"JOIN party_teampoints as ptp ON p.id = ptp.party_id "+
		"WHERE ptp.team_id = ? "+
		"ORDER BY p.created_at DESC",
		team.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parties []domain.Party
	for rows.Next() {
		var party domain.Party
		var date string
		var points float64

		err = rows.Scan(&party.ID, &party.GameID, &date, &points)
		if err != nil {
			return nil, err
		}

		party.TeamPoints = append(party.TeamPoints, domain.TeamPoints{Team: team, Points: points})

		parties = append(parties, party)
	}

	return parties, nil
}

func (s SQLite) PartyCreate(gameID domain.GameID, teamPoints []domain.TeamPoints) (*domain.Party, error) {
	result, err := s.db.Exec("INSERT INTO parties (game_id) VALUES (?)", gameID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var values []any
	var query []string
	for _, teamPoint := range teamPoints {
		query = append(query, "(?,?,?)")
		values = append(values, id, teamPoint.Team.ID, teamPoint.Points)
	}

	_, err = s.db.Exec("INSERT INTO party_teampoints (party_id, team_id, points) VALUES "+strings.Join(query, ","), values...)
	if err != nil {
		return nil, err
	}

	return &domain.Party{ID: domain.PartyID(id), GameID: gameID, TeamPoints: teamPoints}, nil
}

func (s SQLite) GetTableRating() ([]app.TeamRang, error) {
	// TODO: join team data
	rows, err := s.db.Query("SELECT team_id, rating FROM ratings ORDER BY rating DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamRangs []app.TeamRang
	rang := 0

	for rows.Next() {
		rang++

		var team domain.Team
		var rating int

		err = rows.Scan(&team.ID, &rating)
		if err != nil {
			return nil, err
		}

		teamRangs = append(teamRangs, app.TeamRang{Rang: rang, Team: team, Rating: rating})
	}

	return teamRangs, nil
}
