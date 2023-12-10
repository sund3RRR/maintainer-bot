package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func StartHandler(userMessage *tgbotapi.Message, bot *tgbotapi.BotAPI) *tgbotapi.MessageConfig {
	keyboard := GetStartKeyboard()
	msg := tgbotapi.NewMessage(
		userMessage.Chat.ID,
		fmt.Sprintf("Hi! I am %s and my main ability is to notify "+
			"you about new project releases!", bot.Self.UserName),
	)
	msg.ReplyMarkup = keyboard

	return &msg
}

func HandleCommand(command string, userMessage *tgbotapi.Message, bot *tgbotapi.BotAPI) error {
	var msg *tgbotapi.MessageConfig

	switch command {
	case "start":
		msg = StartHandler(userMessage, bot)
	default:
		newMessage := tgbotapi.NewMessage(
			userMessage.Chat.ID,
			fmt.Sprintf("Sorry! Command /%s doesn't exist.", command),
		)
		msg = &newMessage
	}

	_, err := bot.Send(msg)

	return err
}

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

			if command := update.Message.Command(); command != "" {
				logger.Info(
					"Command received",
					zap.String("FromMessage", update.Message.From.UserName),
					zap.String("MessageText", update.Message.Text),
				)
				err := HandleCommand(command, update.Message, bot)
				if err != nil {
					logger.Error(
						"An error occured while sending a message",
						zap.Error(err),
					)
				}
			}
			logger.Info(
				"Message received",
				zap.String("FromMessage", update.Message.From.UserName),
				zap.String("MessageText", update.Message.Text),
			)
		}
	}
}
