package bot

import (
	"errors"
	"strconv"
	"strings"

	"github.com/sund3RRR/maintainer-bot/internal/adapters/db"
	"github.com/sund3RRR/maintainer-bot/internal/adapters/fetcher"

	"gopkg.in/telebot.v3"
)

func parseRepoId(c telebot.Context) (int, error) {
	splitted_query := strings.Split(c.Callback().Data, ":")
	return strconv.Atoi(splitted_query[1])
}

func parseRepoFromUrl(m *telebot.Message) (*db.Repo, error) {
	text, chatID := m.Text, m.Chat.ID
	text = strings.Trim(text, " ")
	text, _ = strings.CutPrefix(text, "http://")
	text, _ = strings.CutPrefix(text, "https://")
	text, _ = strings.CutSuffix(text, ".git")
	text, _ = strings.CutSuffix(text, "/")
	splitted := strings.Split(text, "/")

	if len(splitted) != 3 {
		err := errors.New("Can't parse repo")
		return nil, err
	}

	host, owner, repo := splitted[0], splitted[1], splitted[2]

	newRepo := db.Repo{
		Host:   host,
		Owner:  owner,
		Repo:   repo,
		ChatID: chatID,
	}

	return &newRepo, nil
}

func getRepoFromMessage(message *telebot.Message, f *fetcher.Fetcher) (*db.Repo, error) {
	repo, err := parseRepoFromUrl(message)
	if err != nil {
		return nil, err
	}

	switch repo.Host {
	case "github.com":
		_, err := f.Github.GetGithubRepository(repo)
		if err != nil {
			return nil, err
		}
		tagName, err := f.Github.GetLatestReleaseTagName(repo)
		if err == nil {
			repo.IsRelease = true
			repo.LastTag = tagName

			return repo, nil
		}

		tagName, err = f.Github.GetLatestTagName(repo)
		if err != nil {
			return nil, err
		}

		repo.IsRelease = false
		repo.LastTag = tagName

		return repo, nil
	default:
		return nil, ErrHostIsIncorrect
	}
}

func extractCallbackQuery(callback string) string {
	splitted := strings.Split(callback, ":")

	return splitted[0]
}

func getDefaultSendOptions() *telebot.SendOptions {
	return &telebot.SendOptions{ParseMode: telebot.ModeHTML, ReplyMarkup: &telebot.ReplyMarkup{}}
}
