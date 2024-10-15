package fetcher

import (
	"time"

	"github.com/sund3RRR/maintainer-bot/internal/adapters/db"
	"github.com/sund3RRR/maintainer-bot/internal/config"

	"go.uber.org/zap"
)

type RepoMessage struct {
	ChatID int64
	Title  string
	Text   string
	Link   string
	NewTag string
}

type Fetcher struct {
	RepoUpdatesChan    chan *RepoMessage
	RepoHostingClients *config.RepoHostingClients
	DatabaseService    *db.DatabaseService
	Github             *GithubFetcher
}

func (f *Fetcher) StartRepoFetcher(logger *zap.Logger) {
	for {
		repos, err := f.DatabaseService.GetAllRepos()
		if err != nil {
			logger.Fatal("An error occured while getting all repos", zap.Error(err))
		}
		for _, repo := range *repos {
			var result *RepoMessage

			switch repo.Host {
			case "github.com":
				result = f.Github.FetchRepo(&repo, logger)
			}

			if result != nil && repo.LastTag != result.NewTag {
				f.RepoUpdatesChan <- result

				err := f.DatabaseService.SetLastTagById(result.NewTag, repo.Id)
				if err != nil {
					logger.Error("An error occured while setting last tag by id", zap.Error(err))
				}
			}
		}
		time.Sleep(time.Minute * 5)
	}
}
