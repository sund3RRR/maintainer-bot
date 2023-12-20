package bot

import (
	"app/fetcher"
	"fmt"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func StartBot(token string, repoUpdatesChan chan *fetcher.RepoMessage, logger *zap.Logger) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 3 * time.Second},
	})
	if err != nil {
		logger.Error(
			"An error occured while creating a bot",
			zap.Error(err),
			zap.String("Token", token),
		)
	}
	logger.Info(
		"Bot was successfully created",
		zap.String("Bot UserName", bot.Me.Username),
	)
	storage := memory.NewStorage()
	manager := fsm.NewManager(
		bot,
		nil,
		storage,
		nil,
	)

	logger.Info("The storage was successfully created")

	RegisterHandlers(manager, logger, bot)

	go StartRepoSender(repoUpdatesChan, bot, logger)

	bot.Start()
}

func StartRepoSender(c chan *fetcher.RepoMessage, bot *telebot.Bot, logger *zap.Logger) {
	logger.Info("Repo sender goroutine was successfully created")
	for repo := range c {
		logger.Info("New send message request received, processing")

		if repo.Text != "" {
			repo.Text += "\n\n"
		}
		message := fmt.Sprintf("%s\n\n%s%s", repo.Title, repo.Text, repo.Link)

		_, err := bot.Send(&telebot.Chat{ID: int64(repo.ChatID)}, message, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
		if err != nil {
			logger.Error(
				"An error occured while sending repo message",
				zap.Error(err),
				zap.String("ChatID", fmt.Sprintf("%d", repo.ChatID)),
				zap.String("Repo message", message),
			)
		}
		logger.Info("Successully send repo message")
	}
	logger.Info("WTF")
}
