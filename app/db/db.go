package db

import (
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	Id        int    `db:"id"`
	ChatID    int    `db:"chat_id"`
	Host      string `db:"host"`
	Owner     string `db:"owner"`
	Repo      string `db:"repo"`
	LastTag   string `db:"last_tag"`
	IsRelease bool   `db:"is_release"`
}

type User struct {
	Id     int `db:"id"`
	ChatID int `db:"chat_id"`
}

var DBInstance *sqlx.DB

func PrepareDb(conn *sqlx.DB) error {
	_, err := conn.Exec("DROP TABLE IF EXISTS repos;")

	_, err = conn.Exec(
		`CREATE TABLE IF NOT EXISTS repos(
            id SERIAL PRIMARY KEY,
			chat_id INTEGER,
            host TEXT,
            owner TEXT,
            repo TEXT,
            last_tag TEXT,
            is_release BOOLEAN
		);`,
	)
	return err
}
