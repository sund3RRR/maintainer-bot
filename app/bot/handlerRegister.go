package bot

import (
	"app/bot/handlers"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func RegisterHandlers(manager *fsm.Manager, logger *zap.Logger, bot *telebot.Bot) {
	manager.Bind("/start", fsm.AnyState, func(c telebot.Context, state fsm.Context) error {
		err := handlers.StartHandler(c, state, logger, bot)
		if err != nil {
			logger.Error(
				"An error occured while handling /start",
				zap.Error(err),
			)
		}
		return err
	})

	manager.Bind("/add_repo", fsm.DefaultState, func(c telebot.Context, state fsm.Context) error {
		err := handlers.AddRepoHandler(c, state, logger, bot)
		if err != nil {
			logger.Error(
				"An error occured while handling /add_repo",
				zap.Error(err),
			)
		}
		return err
	})

	manager.Handle(fsm.F(telebot.OnText, handlers.AddRepoState), func(c telebot.Context, state fsm.Context) error {
		err := handlers.OnRepoEntered(c, state, logger, bot)
		if err != nil {
			logger.Error(
				"An error occured while handling /add_repo",
				zap.Error(err),
			)
		}
		return err
	})
}
