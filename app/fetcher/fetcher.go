package fetcher

import (
	"app/config"
	repodb "app/db"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type RepoMessage struct {
	ChatID int
	Title  string
	Text   string
	Link   string
	NewTag string
}

var RepoHostingClientsVar *config.RepoHostingClients

func StartRepoFetcher(c chan *RepoMessage, db *sqlx.DB, repoHostingClients *config.RepoHostingClients, logger *zap.Logger) {
	RepoHostingClientsVar = repoHostingClients

	for {
		repos := []repodb.Repo{}
		err := db.Select(&repos, `SELECT * FROM repos`)
		if err != nil {
			logger.Error(
				"An error occured while select repos",
				zap.Error(err),
			)
		}
		for _, repo := range repos {
			logger.Info("New repo fetching...", zap.String("Repo", repo.Repo))
			var result *RepoMessage
			switch repo.Host {
			case "github.com":
				result = FetchGithubRepo(&repo, repoHostingClients.GitHub, logger)
			}
			if result != nil {
				logger.Info("Result is not nil")
				c <- result

				_, err := db.Exec("UPDATE repos SET last_tag = $1 WHERE id = $2", result.NewTag, repo.Id)
				if err != nil {
					logger.Error(
						"An error occured while updating repos last_tag",
						zap.Error(err),
					)
				}
			}
		}
		time.Sleep(time.Minute)
	}
}
