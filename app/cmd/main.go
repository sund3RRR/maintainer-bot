package main

import (
	"app/bot"
	appConfig "app/config"
	appDB "app/db"
	"app/fetcher"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v57/github"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func main() {
	config, err := appConfig.NewConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := config.ZapConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
	)
	db, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	appDB.DBInstance = db

	logger.Info(
		"Successfully connected to PostgeSQL",
		zap.String("Host", config.Postgres.Host),
		zap.Int("Port", config.Postgres.Port),
		zap.String("User", config.Postgres.User),
		zap.String("Database", config.Postgres.Database),
	)

	if len(os.Args) > 1 && os.Args[1] == "--prepare-db" {
		appDB.PrepareDb(db)
		logger.Info("Prepare database request complete successfully")
	}

	logger.Info("Starting telegram bot...")

	githubClient := github.NewClient(nil).WithAuthToken(config.RepoHostingApis.GithubToken)

	repoHostingClients := &appConfig.RepoHostingClients{
		GitHub: githubClient,
	}
	repoUpdatesChan := make(chan *fetcher.RepoMessage)

	go fetcher.StartRepoFetcher(repoUpdatesChan, db, repoHostingClients, logger)

	bot.StartBot(config.TelegramBot.Token, repoUpdatesChan, logger)
}
