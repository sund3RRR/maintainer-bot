package handlers

import (
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func GetStartKeyboard() *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				telebot.ReplyButton{
					Text: "/add_repo",
				},
				telebot.ReplyButton{
					Text: "/remove_repo",
				},
			},
			{
				telebot.ReplyButton{
					Text: "/list_repos",
				},
			},
		},
	}
	return keyboard
}

func StartHandler(c telebot.Context, state fsm.Context, logger *zap.Logger, bot *telebot.Bot) error {
	logger.Info(
		"Received /start, handling a command...",
		zap.String("Sender username", c.Sender().Username),
	)
	keyboard := GetStartKeyboard()

	if state != nil {
		err := state.Finish(true)
		if err != nil {
			logger.Error(
				"An error occured while finishing state",
				zap.Error(err),
				zap.String("Context message", c.Text()),
				zap.String("Context data", c.Data()),
			)
		}
	}

	if err := c.Send(fmt.Sprintf("Hi! I am %s and my main ability is to notify "+
		"you about new project releases!", bot.Me.FirstName)); err != nil {
		return err
	}

	if err := c.Send("How can I help you?", keyboard); err != nil {
		return err
	}

	return nil
}
