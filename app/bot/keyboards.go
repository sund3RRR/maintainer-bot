package bot

import (
	"app/db"
	"fmt"

	"gopkg.in/telebot.v3"
)

func getHomeKeyboard() *telebot.ReplyMarkup {
	keyboard := telebot.ReplyMarkup{
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				telebot.ReplyButton{
					Text: "/home",
				},
			},
		},
	}
	keyboard.ResizeKeyboard = true
	return &keyboard
}

func getStartKeyboard() *telebot.ReplyMarkup {
	keyboard := telebot.ReplyMarkup{
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
	keyboard.ResizeKeyboard = true
	return &keyboard
}

func getRemoveRepoKeyboard(repos *[]db.Repo) *telebot.ReplyMarkup {
	keyboard := telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{}}}

	for _, repo := range *repos {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []telebot.InlineButton{
			{
				Text: fmt.Sprintf("%s: %s/%s", repo.Host, repo.Owner, repo.Repo),
				Data: "remove_repo:" + fmt.Sprint(repo.Id),
			},
		})
	}

	return &keyboard
}
