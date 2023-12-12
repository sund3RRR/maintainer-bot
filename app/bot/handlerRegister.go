package bot

import (
	"app/bot/handlers"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	tgbot "gopkg.in/telebot.v3"
)

func RegisterHandlers(manager *fsm.Manager, logger *zap.Logger, bot *tgbot.Bot) {
	manager.Bind("/start", fsm.AnyState, func(c tgbot.Context, state fsm.Context) error {
		err := handlers.StartHandler(c, state, logger, bot)
		if err != nil {
			logger.Error(
				"An error occured while handling /start",
				zap.Error(err),
			)
		}
		return err
	})
}
