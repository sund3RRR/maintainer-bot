package main

import (
	"app/bot"
	"app/config"
	"app/db"
	"app/fetcher"
	"log"
	"sync"

	"github.com/google/go-github/v57/github"
	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.NewConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := cfg.ZapConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()

	dbService := db.DatabaseService{}
	err = dbService.Connect(cfg)
	if err != nil {
		logger.Fatal(
			"An error occured while trying to connect postgreSQL",
			zap.Error(err),
			zap.String("Host", cfg.Postgres.Host),
			zap.Int("Port", cfg.Postgres.Port),
			zap.String("User", cfg.Postgres.User),
			zap.String("Database", cfg.Postgres.Database),
		)
	}

	defer dbService.DB.Close()
	logger.Info("Successfully connected to PostgeSQL")

	err = dbService.PrepareDb()
	if err != nil {
		logger.Fatal("An error occured while trying to prepare DB", zap.Error(err))
	}

	logger.Info("Starting telegram bot...")

	githubClient := github.NewClient(nil).WithAuthToken(cfg.RepoHostingApis.GithubToken)

	repoUpdatesChan := make(chan *fetcher.RepoMessage)

	var wg sync.WaitGroup
	wg.Add(3)

	f := fetcher.Fetcher{
		RepoUpdatesChan: repoUpdatesChan,
		DatabaseService: &dbService,
		Github:          &fetcher.GithubFetcher{Client: githubClient},
	}

	go func() {
		defer wg.Done()
		f.StartRepoFetcher(logger)
	}()

	botService := bot.BotService{
		Token:           cfg.TelegramBot.Token,
		RepoUpdatesChan: repoUpdatesChan,
		DatabaseService: &dbService,
		Logger:          logger,
		Fetcher:         &f,
	}

	go func() {
		defer wg.Done()
		botService.StartBot()
	}()

	go func() {
		defer wg.Done()
		botService.StartRepoSender()
	}()

	wg.Wait()
}
