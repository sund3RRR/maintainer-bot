package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func StartBot(token string, logger *zap.Logger) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error(
			"An error occured while creating a bot",
			zap.Error(err),
			zap.String("Token", token),
		)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	logger.Info("Authorized on account", zap.String("Account", bot.Self.UserName))

	for update := range updates {
		if update.Message != nil {
			logger.Info(
				"Message received",
				zap.String("FromMessage", update.Message.From.UserName),
				zap.String("MessageText", update.Message.Text),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				logger.Error(
					"An error occured while sending a message",
					zap.Error(err),
				)
			}
		}
	}
}
