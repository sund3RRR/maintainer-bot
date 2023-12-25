package bot

import (
	"github.com/vitaliy-ukiru/fsm-telebot"
	"gopkg.in/telebot.v3"
)

// Registers all handlers that are used by this bot including fsm handlers
func (botService *BotService) RegisterAllHandlers() {
	botService.manager.Bind("/start", fsm.AnyState, botService.StartHandler)
	botService.manager.Bind("/home", fsm.AnyState, botService.HomeHandler)

	// Add Repo
	botService.manager.Bind("/add_repo", fsm.DefaultState, botService.AddRepoHandler)
	botService.manager.Handle(fsm.F(telebot.OnText, AddRepoState), botService.OnRepoEntered)
	// Add Repo

	botService.manager.Bind("/remove_repo", fsm.DefaultState, botService.RemoveRepoHandler)

	// List Repo
	botService.manager.Bind("/list_repos", fsm.DefaultState, botService.ListRepoHandler)
	// List Repo

	botService.manager.Handle(fsm.F(telebot.OnCallback, fsm.AnyState), func(c telebot.Context, state fsm.Context) error {
		switch extractCallbackQuery(c.Callback().Data) {
		case "remove_repo":
			err := botService.OnRemoveRepo(c, state)
			return err
		default:
			return state.Finish(c.Data() != "")
		}
	})
}
