package db

import (
	"app/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseService struct {
	DB *sqlx.DB
}

type Repo struct {
	Id        int    `db:"id"`
	ChatID    int64  `db:"chat_id"`
	Host      string `db:"host"`
	Owner     string `db:"owner"`
	Repo      string `db:"repo"`
	LastTag   string `db:"last_tag"`
	IsRelease bool   `db:"is_release"`
}

var DBInstance *sqlx.DB

func (dbService *DatabaseService) PrepareDb() error {
	_, err := dbService.DB.Exec(
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

func (dbService *DatabaseService) GetDatabaseUrl(cfg *config.AppConfig) string {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	return databaseUrl
}

func (dbService *DatabaseService) Connect(cfg *config.AppConfig) error {
	db, err := sqlx.Connect("postgres", dbService.GetDatabaseUrl(cfg))
	if err != nil {
		return err
	}

	dbService.DB = db

	return nil
}

func (dbService *DatabaseService) GetAllRepos() (*[]Repo, error) {
	var repos []Repo

	query := "SELECT * FROM repos;"
	err := dbService.DB.Select(&repos, query)

	return &repos, err
}

func (dbService *DatabaseService) DeleteRepoWhereId(id int) error {
	_, err := dbService.DB.Exec("DELETE FROM repos WHERE id = $1;", id)
	return err
}
func (dbService *DatabaseService) SetLastTagById(newTag string, id int) error {
	query := "UPDATE repos SET last_tag = $1 WHERE id = $2"
	_, err := dbService.DB.Exec(query, newTag, id)
	return err
}

func (dbService *DatabaseService) GetRepoWhereId(id int) (*Repo, error) {
	var repo Repo

	query := "SELECT * FROM repos WHERE id = $1;"
	err := dbService.DB.Get(&repo, query, id)

	return &repo, err
}

func (dbService *DatabaseService) GetReposWhereChatId(chat_id int64) (*[]Repo, error) {
	var repos []Repo
	err := dbService.DB.Select(&repos, "SELECT * FROM repos WHERE chat_id = $1 ORDER BY host;", chat_id)

	return &repos, err
}

func (dbService *DatabaseService) IsRepoAlreadyExist(repo *Repo) (bool, error) {
	var count int

	query := "SELECT COUNT(*) FROM repos WHERE chat_id=$1 AND host=$2 AND owner=$3 AND repo=$4"
	err := dbService.DB.Get(&count, query, repo.ChatID, repo.Host, repo.Owner, repo.Repo)

	return count != 0, err
}

func (dbService *DatabaseService) AddRepo(repo *Repo) error {
	query := `INSERT INTO repos (host, owner, repo, chat_id, last_tag, is_release)
	VALUES (:host, :owner, :repo, :chat_id, :last_tag, :is_release);`
	_, err := dbService.DB.NamedExec(query, repo)

	return err
}
