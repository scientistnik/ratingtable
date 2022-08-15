-- +migrate Up
-- TODO add autoincrement and uniq
CREATE TABLE IF NOT EXISTS games (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT,
  type INTEGER
);

CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY,
  links JSON
);

CREATE TABLE IF NOT EXISTS teams (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT,
  game_id REFERENCES games
);

CREATE TABLE IF NOT EXISTS user_team (
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER REFERENCES users,
  team_id INTEGER REFERENCES teams
);

CREATE TABLE IF NOT EXISTS parties (
  id INTEGER NOT NULL PRIMARY KEY,
  game_id INTEGER REFERENCES games,
  created_at DATETIME
);

CREATE TABLE IF NOT EXISTS party_teampoints (
  id INTEGER NOT NULL PRIMARY KEY,
  party_id INTEGER REFERENCES parties,
  team_id INTEGER REFERENCES teams,
  points FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS rating_changes (
  id INTEGER NOT NULL PRIMARY KEY,
  team_id INTEGER REFERENCES teams,
  party_id INTEGER REFERENCES parties,
  points_change INTEGER NOT NULL,
  updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS ratings (
  id INTEGER NOT NULL PRIMARY KEY,
  team_id INTEGER REFERENCES teams,
  rating INTEGER NOT NULL,
  updated_at DATETIME
);

-- +migrate Down
DROP TABLE games;
DROP TABLE users;
DROP TABLE teams;
DROP TABLE user_team;
DROP TABLE parties;
DROP TABLE party_team;
DROP TABLE rating_changes;
DROP TABLE ratings;