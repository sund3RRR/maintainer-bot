package handlers

import (
	appDB "app/db"
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func ListRepoHandler(c telebot.Context, state fsm.Context, logger *zap.Logger, bot *telebot.Bot) error {
	db := appDB.DBInstance

	var reposInfo []RepoInfo
	err := db.Select(&reposInfo, "SELECT host, owner, repo FROM repos WHERE chat_id=$1 ORDER BY host;", c.Chat().ID)
	if err != nil {
		logger.Error(
			"An error occured while trying to SELECT repoInfo",
			zap.Error(err),
		)
	}

	if len(reposInfo) == 0 {
		message_text := "Hmmm, it looks like you don't have any repositories yet. Wanna add a couple?"
		err = c.Send(message_text, &telebot.SendOptions{ReplyMarkup: GetStartKeyboard()})
		return err
	}

	resultStr := ""
	for i, repo := range reposInfo {
		repoLink := fmt.Sprintf(`<a href="https://%s/%s/%s/">link</a>`, repo.Host, repo.Owner, repo.Repo)
		resultStr += fmt.Sprintf("%d) <code>%s:%s/%s</code> %s\n", i+1, repo.Host, repo.Owner, repo.Repo, repoLink)
	}

	err = c.Send(resultStr, &telebot.SendOptions{ParseMode: telebot.ModeHTML, ReplyMarkup: GetStartKeyboard(), DisableWebPagePreview: true})

	return err
}
