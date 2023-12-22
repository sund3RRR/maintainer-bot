package handlers

import (
	"app/db"
	"app/fetcher"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

var ErrHostIsIncorrect = errors.New("Repo host is incorrrect")

const (
	AddRepoState fsm.State = "add_repo"
)

type Repo struct {
	ID        int    `db:"id"`
	Host      string `db:"host"`
	Owner     string `db:"owner"`
	Repo      string `db:"repo"`
	ChatID    int64  `db:"chat_id"`
	LastTag   string `db:"last_tag"`
	IsRelease bool   `db:"is_release"`
}
type RepoInfo struct {
	Host  string `db:"host"`
	Owner string `db:"owner"`
	Repo  string `db:"repo"`
}

func AddRepoHandler(c telebot.Context, state fsm.Context, logger *zap.Logger, bot *telebot.Bot) error {
	logger.Info(
		"Received /add_repo, handling command...",
		zap.String("Sender username", c.Sender().Username),
	)

	err := state.Set(AddRepoState)
	if err != nil {
		logger.Error(
			fmt.Sprintf("An error occured while setting a state %s", AddRepoState),
			zap.Error(err),
			zap.String("Sender username", c.Sender().Username),
		)
	}

	err = c.Send("No problem! Just send me a link to the repository (GitHub only for now)", GetHomeKeyboard())

	return err
}

func OnRepoEntered(c telebot.Context, state fsm.Context, logger *zap.Logger, bot *telebot.Bot) error {
	var githubErrResponse *github.ErrorResponse

	repo, err := parseMessage(c.Message())
	if err != nil {
		messageText := "Sorry, but I can't add the repository :(" +
			"\n Please enter the repo in the format 'https://<u>host</u>/<u>owner</u>/<u>repo</u>'"
		err := c.Send(messageText, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
		if err != nil {
			logger.Error(
				"An error occured while trying to send 'formating repo error' for user",
				zap.Error(err),
			)
		}
		return nil
	} else if err := checkRepoIsValid(repo); err != nil {
		if errors.Is(err, ErrHostIsIncorrect) {
			messageText := fmt.Sprintf("Host <code>%s</code> is incorrect or doesn't support yet :(", repo.Host)
			err := c.Send(messageText, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
			if err != nil {
				logger.Error(
					"An error occured while trying to send 'host is incorrect' for user",
					zap.Error(err),
				)
			}
			return nil
		} else if errors.As(err, &githubErrResponse) {
			messageText := fmt.Sprintf("Repo <code>%s</code> is incorrect/private or doesn't exist :(", repo.Repo)
			err := c.Send(messageText, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
			if err != nil {
				logger.Error(
					"An error occured while trying to send 'host is incorrect' for user",
					zap.Error(err),
				)
			}

			return nil
		}
	}

	var count int
	err = db.DBInstance.Get(
		&count,
		"SELECT COUNT(*) FROM repos WHERE chat_id=$1 AND host=$2 AND owner=$3 AND repo=$4",
		repo.ChatID,
		repo.Host,
		repo.Owner,
		repo.Repo,
	)
	if err != nil {
		logger.Error(
			"An error occured while trying to db.Get repo",
			zap.Error(err),
			zap.String("Repo owner", repo.Owner),
			zap.String("Repo", repo.Repo),
		)
	}

	if count > 0 {
		messageText := fmt.Sprintf("<code>%s:%s/%s</code> is already exist", repo.Host, repo.Owner, repo.Repo)
		err := c.Send(messageText, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
		if err != nil {
			logger.Error(
				"An error occured while trying to send 'repo is already exist' for user",
				zap.Error(err),
			)
		}
		return nil
	}

	tagName, isRelease := checkRepoHasReleases(repo)
	if !isRelease {
		tagName = getRepoLastTagName(repo, logger)
	}

	repo.LastTag = tagName
	repo.IsRelease = isRelease

	_, err = db.DBInstance.NamedExec(
		`INSERT INTO repos (host, owner, repo, chat_id, last_tag, is_release)
		VALUES (:host, :owner, :repo, :chat_id, :last_tag, :is_release);`,
		repo,
	)
	if err != nil {
		logger.Error(
			"An error occured while trying to db.NamedExec new repo",
			zap.Error(err),
		)
	}

	err = state.Finish(c.Data() != "")
	if err != nil {
		logger.Error(
			"An error occured while trying to finish AddRepo state",
			zap.Error(err),
		)
	}

	messageText := fmt.Sprintf("Repo <code>%s:%s/%s</code> successfully added!", repo.Host, repo.Owner, repo.Repo)
	err = c.Send(messageText, &telebot.SendOptions{ParseMode: telebot.ModeHTML, ReplyMarkup: GetStartKeyboard()})
	if err != nil {
		logger.Error(
			"An error occured while sending message to user",
			zap.Error(err),
		)
	}

	return err
}

func parseMessage(m *telebot.Message) (*Repo, error) {
	text, chatID := m.Text, m.Chat.ID
	text = strings.Trim(text, " ")
	text, _ = strings.CutPrefix(text, "http://")
	text, _ = strings.CutPrefix(text, "https://")
	text, _ = strings.CutSuffix(text, ".git")
	splitted := strings.Split(text, "/")

	if len(splitted) != 3 {
		err := errors.New("Can't parse repo")
		return nil, err
	}

	host, owner, repo := splitted[0], splitted[1], splitted[2]

	newRepo := Repo{
		Host:   host,
		Owner:  owner,
		Repo:   repo,
		ChatID: chatID,
	}

	return &newRepo, nil
}

func checkRepoIsValid(repo *Repo) error {
	switch repo.Host {
	case "github.com":
		githubClient := fetcher.RepoHostingClientsVar.GitHub

		_, _, err := githubClient.Repositories.Get(context.Background(), repo.Owner, repo.Repo)

		return err
	default:
		return ErrHostIsIncorrect
	}
}

func checkRepoHasReleases(repo *Repo) (string, bool) {
	tagName := ""
	switch repo.Host {
	case "github.com":
		githubClient := fetcher.RepoHostingClientsVar.GitHub

		lastRelease, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), repo.Owner, repo.Repo)
		if err != nil {
			return "", false
		}
		tagName = lastRelease.GetTagName()
	default:
		return "", false
	}

	return tagName, true
}

func getRepoLastTagName(repo *Repo, logger *zap.Logger) string {
	tagName := ""
	switch repo.Host {
	case "github.com":
		githubClient := fetcher.RepoHostingClientsVar.GitHub

		tags, _, err := githubClient.Repositories.ListTags(context.Background(), repo.Owner, repo.Repo, &github.ListOptions{Page: 1})
		if err != nil {
			logger.Error(
				"An error occured while getting tags from github repo",
				zap.Error(err),
				zap.String("Repo owner", repo.Owner),
				zap.String("Repo", repo.Repo),
			)
		}
		tagName = tags[0].GetName()
	default:
		return ""
	}

	return tagName
}
