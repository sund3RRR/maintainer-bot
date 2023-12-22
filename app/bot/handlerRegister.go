package bot

import (
	"app/bot/handlers"
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

// Registers all handlers that are used by this bot including fsm handlers
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
	manager.Bind("/home", fsm.AnyState, func(c telebot.Context, state fsm.Context) error {
		err := state.Finish(c.Data() != "")
		if err != nil {
			logger.Error(
				"An error occured while trying to finish state",
				zap.Error(err),
			)
		}

		err = c.Send("How can I help you?", &telebot.SendOptions{ReplyMarkup: handlers.GetStartKeyboard()})

		return err
	})

	//
	// Add Repo
	//
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
	//
	// Add Repo
	//

	//
	// Remove Repo
	//
	manager.Bind("/remove_repo", fsm.DefaultState, func(c telebot.Context, state fsm.Context) error {
		err := handlers.RemoveRepoHandler(c, state, logger, bot)
		if err != nil {
			logger.Error(
				"An error occured while handling /remove_repo",
				zap.Error(err),
			)
		}
		return err
	})
	manager.Handle(fsm.F(telebot.OnCallback, fsm.AnyState), func(c telebot.Context, state fsm.Context) error {
		switch ExtractCallbackQuery(c.Callback().Data) {
		case "remove_repo":
			return handlers.OnRemoveRepoEntered(c, state, logger, bot)
		default:
			return state.Finish(c.Data() != "")
		}
	})
	//
	// Remove Repo
	//
}

func ExtractCallbackQuery(callback string) string {
	splitted := strings.Split(callback, ":")

	return splitted[0]
}
