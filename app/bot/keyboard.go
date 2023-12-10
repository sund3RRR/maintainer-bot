package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetStartKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Add repo"),
			tgbotapi.NewKeyboardButton("Remove repo"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("List repos"),
		),
	)
	return keyboard
}
