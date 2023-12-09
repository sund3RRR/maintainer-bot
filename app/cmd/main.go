package main

import (
	"app/bot"
	"app/config"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func main() {
	fmt.Print()
	config, err := config.NewConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := config.ZapConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
	)
	fmt.Println(databaseUrl)
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		logger.Error("Unable connect to database", zap.Error(err))
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	logger.Info(
		"Successfully connected to PostgeSQL",
		zap.String("Host", config.Postgres.Host),
		zap.Int("Port", config.Postgres.Port),
		zap.String("User", config.Postgres.User),
		zap.String("Database", config.Postgres.Database),
	)

	logger.Info("Starting telegram bot...")

	bot.StartBot(config.TelegramBot.Token, logger)
}
