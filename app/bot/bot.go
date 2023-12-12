package bot

import (
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"go.uber.org/zap"
	tgbot "gopkg.in/telebot.v3"
)

func StartBot(token string, logger *zap.Logger) {
	bot, err := tgbot.NewBot(tgbot.Settings{
		Token:  token,
		Poller: &tgbot.LongPoller{Timeout: 3 * time.Second},
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
	logger.Info(
		"The storage was successfully created",
	)
	RegisterHandlers(manager, logger, bot)

	bot.Start()
}
