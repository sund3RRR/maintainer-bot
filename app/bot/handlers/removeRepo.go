package handlers

import (
	dbApp "app/db"
	"fmt"
	"strconv"
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

const (
	RemoveRepoState fsm.State = "add_repo"
)

func getRemoveRepoKeyboard(repos []*Repo) *telebot.ReplyMarkup {
	keyboard := telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{}}}

	for _, repo := range repos {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []telebot.InlineButton{
			{
				Text: fmt.Sprintf("%s: %s/%s", repo.Host, repo.Owner, repo.Repo),
				Data: "remove_repo:" + fmt.Sprint(repo.ID),
			},
		})
	}

	return &keyboard
}
func RemoveRepoHandler(c telebot.Context, state fsm.Context, logger *zap.Logger, bot *telebot.Bot) error {
	logger.Info(
		"Received /remove_repo, handling command...",
		zap.String("Sender username", c.Sender().Username),
	)

	db := dbApp.DBInstance

	var repos []*Repo
	err := db.Select(&repos, "SELECT * FROM repos WHERE chat_id = $1;", c.Chat().ID)
	if err != nil {
		logger.Error(
			"An error occured while trying to SELECT repos in RemoveRepoHandler",
			zap.Error(err),
		)
	}

	if len(repos) == 0 {
		message_text := "Hmmm, it looks like you don't have any repositories yet. Wanna add a couple?"
		err = c.Send(message_text, &telebot.SendOptions{ReplyMarkup: GetStartKeyboard()})
		return err
	}

	err = state.Set(RemoveRepoState)
	if err != nil {
		logger.Error(
			fmt.Sprintf("An error occured while setting a state %s", RemoveRepoState),
			zap.Error(err),
			zap.String("Sender username", c.Sender().Username),
		)
	}
	err = c.Send("No problem, man!", &telebot.SendOptions{ReplyMarkup: GetHomeKeyboard()})
	if err != nil {
		logger.Error(
			"An error occured while trying to send home keyboard",
			zap.Error(err),
		)
	}
	keyboard := getRemoveRepoKeyboard(repos)
	err = c.Send("Select the repository you want to remove", &telebot.SendOptions{ReplyMarkup: keyboard})

	return err
}

func OnRemoveRepoEntered(c telebot.Context, state fsm.Context, logger *zap.Logger, bot *telebot.Bot) error {
	splitted_query := strings.Split(c.Callback().Data, ":")
	id, err := strconv.Atoi(splitted_query[1])
	if err != nil {
		logger.Error(
			"An error occured while trying to parse ID from callback query",
			zap.Error(err),
			zap.String("Callback", splitted_query[0]),
			zap.String("ID", splitted_query[1]),
		)
		return err
	}

	db := dbApp.DBInstance

	type RemoveRepoDB struct {
		Host  string `db:"host"`
		Owner string `db:"owner"`
		Repo  string `db:"repo"`
	}
	var removeRepoDB RemoveRepoDB
	err = db.Get(&removeRepoDB, "SELECT host, owner, repo FROM repos WHERE id = $1;", id)
	if err != nil {
		logger.Error(
			"An error occured while trying to GET repo",
			zap.Error(err),
		)
	}
	_, err = db.Exec("DELETE FROM repos WHERE id = $1;", id)
	if err != nil {
		logger.Error(
			"An error occured while trying to exec DELETE repo from repos",
			zap.Error(err),
		)
		return err
	}

	err = state.Finish(c.Data() != "")
	if err != nil {
		logger.Error(
			"An error occured while trying to finish RemoveRepoState ",
			zap.Error(err),
		)
	}

	_, err = bot.Edit(c.Callback(), &telebot.ReplyMarkup{})
	if err != nil {
		logger.Error(
			"An error occured while trying to remove inline keyboard markup",
			zap.Error(err),
		)
	}
	message_text := fmt.Sprintf("Repo <code>%s:%s/%s</code> successfully removed", removeRepoDB.Host, removeRepoDB.Owner, removeRepoDB.Repo)
	err = c.Send(message_text, &telebot.SendOptions{ParseMode: telebot.ModeHTML, ReplyMarkup: GetStartKeyboard()})

	return err
}
