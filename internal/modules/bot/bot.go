package bot

import (
	"fmt"
	"time"

	"github.com/sund3RRR/maintainer-bot/internal/adapters/db"
	"github.com/sund3RRR/maintainer-bot/internal/adapters/fetcher"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type BotService struct {
	Token           string
	RepoUpdatesChan chan *fetcher.RepoMessage
	DatabaseService *db.DatabaseService
	Logger          *zap.Logger
	Fetcher         *fetcher.Fetcher
	manager         *fsm.Manager
	bot             *telebot.Bot
}

func (botService *BotService) StartBot() {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botService.Token,
		Poller: &telebot.LongPoller{Timeout: 3 * time.Second},
	})

	botService.bot = bot

	if err != nil {
		botService.Logger.Error(
			"An error occured while creating a bot",
			zap.Error(err),
			zap.String("Token", botService.Token),
		)
	}

	botService.Logger.Info(
		"Bot was successfully created",
		zap.String("Bot UserName", bot.Me.Username),
	)

	storage := memory.NewStorage()
	botService.manager = fsm.NewManager(
		bot,
		nil,
		storage,
		nil,
	)

	botService.Logger.Info("The storage was successfully created")

	botService.RegisterAllHandlers()

	botService.bot.Start()
}

func (botService *BotService) StartRepoSender() {
	botService.Logger.Info("Repo sender goroutine was successfully created")

	for repo := range botService.RepoUpdatesChan {
		botService.Logger.Info("New send message request received, processing")

		if repo.Text != "" {
			repo.Text += "\n"
		}
		message := fmt.Sprintf("%s\n\n%s%s", repo.Title, repo.Text, repo.Link)

		_, err := botService.bot.Send(&telebot.Chat{ID: repo.ChatID}, message, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
		if err != nil {
			botService.Logger.Error(
				"An error occured while sending repo message",
				zap.Error(err),
				zap.String("ChatID", fmt.Sprintf("%d", repo.ChatID)),
				zap.String("Repo message", message),
			)
		}
		botService.Logger.Info("Successully send repo message")
	}
}
