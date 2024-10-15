package bot

import (
	"errors"
	"fmt"

	"github.com/sund3RRR/maintainer-bot/internal/adapters/fetcher"

	"github.com/google/go-github/v57/github"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

// This is a handler on command `/list_repos` and it will send the user a list of all the repositories that he has added in format
// `№) <host>:<owner>/<repo> v0.0.0 <repo_link>`
func (botService *BotService) ListRepoHandler(c telebot.Context, state fsm.Context) error {
	repos, err := botService.DatabaseService.GetReposWhereChatId(c.Chat().ID)
	if err != nil {
		botService.Logger.Error(
			"An error occured while trying to SELECT repoInfo",
			zap.Error(err),
		)
	}

	if len(*repos) == 0 {
		message_text := "Hmmm, it looks like you don't have any repositories yet. Wanna add a couple?"
		err = c.Send(message_text, &telebot.SendOptions{ReplyMarkup: getStartKeyboard()})
		return err
	}

	resultStr := ""
	for i, repo := range *repos {
		repoLink := fmt.Sprintf(`<a href="https://%s/%s/%s/">link</a>`, repo.Host, repo.Owner, repo.Repo)
		resultStr += fmt.Sprintf("%d) <code>%s:%s/%s %s</code> %s\n", i+1, repo.Host, repo.Owner, repo.Repo, repo.LastTag, repoLink)
	}

	sendOptions := &telebot.SendOptions{
		ParseMode:             telebot.ModeHTML,
		ReplyMarkup:           getStartKeyboard(),
		DisableWebPagePreview: true,
	}

	err = c.Send(resultStr, sendOptions)

	return err
}

// This is the default handler on command `/start`, that represents the bot, sends start keyboard and resets all states
func (botService *BotService) StartHandler(c telebot.Context, state fsm.Context) error {
	botService.Logger.Info(
		"Received /start, handling a command...",
		zap.String("Sender username", c.Sender().Username),
	)
	keyboard := getStartKeyboard()

	if state != nil {
		err := state.Finish(true)
		if err != nil {
			botService.Logger.Error(
				"An error occured while finishing state",
				zap.Error(err),
				zap.String("Context message", c.Text()),
				zap.String("Context data", c.Data()),
			)
		}
	}

	err := c.Send(fmt.Sprintf("Hi! I am %s and my main ability is to notify "+
		"you about new project releases!", botService.bot.Me.FirstName))
	if err != nil {
		return err
	}

	err = c.Send("How can I help you?", &telebot.SendOptions{ReplyMarkup: keyboard})

	return err
}

func (botService *BotService) HomeHandler(c telebot.Context, state fsm.Context) error {
	err := state.Finish(c.Data() != "")
	if err != nil {
		botService.Logger.Error("An error occured while trying to finish state", zap.Error(err))
	}

	sendOptions := getDefaultSendOptions()
	sendOptions.ReplyMarkup = getStartKeyboard()

	err = c.Send("How can I help you?", sendOptions)

	return err
}

// This is a handler on command `/remove_repo“, that will start remove repo procedure: send inline keyboard
// with user's repositories on buttons, start remove_repo state and hide start keyboard
func (botService *BotService) RemoveRepoHandler(c telebot.Context, state fsm.Context) error {
	botService.Logger.Info(
		"Received /remove_repo, handling command...",
		zap.String("Sender username", c.Sender().Username),
	)

	repos, err := botService.DatabaseService.GetReposWhereChatId(c.Chat().ID)
	if err != nil {
		botService.Logger.Error("An error occured while trying to SELECT repos in RemoveRepoHandler", zap.Error(err))
	}

	if len(*repos) == 0 {
		message_text := "Hmmm, it looks like you don't have any repositories yet. Wanna add a couple?"
		err = c.Send(message_text, &telebot.SendOptions{ReplyMarkup: getStartKeyboard()})
		return err
	}

	err = state.Set(RemoveRepoState)
	if err != nil {
		botService.Logger.Error(
			fmt.Sprintf("An error occured while setting a state %s", RemoveRepoState),
			zap.Error(err),
			zap.String("Sender username", c.Sender().Username),
		)
	}
	err = c.Send("No problem, man!", &telebot.SendOptions{ReplyMarkup: getHomeKeyboard()})
	if err != nil {
		botService.Logger.Error(
			"An error occured while trying to send home keyboard",
			zap.Error(err),
		)
	}
	keyboard := getRemoveRepoKeyboard(repos)
	err = c.Send("Select the repository you want to remove", &telebot.SendOptions{ReplyMarkup: keyboard})

	return err
}

// This is a callback on RemoveRepoHandler, that will remove repo from DB, hide inline keyboard, finish remove_repo state
// and send start keyboard
func (botService *BotService) OnRemoveRepo(c telebot.Context, state fsm.Context) error {
	id, err := parseRepoId(c)
	if err != nil {
		botService.Logger.Error("An error occured while trying to parse repo ID", zap.Error(err))
		return c.Send("Sorry! An internal error occured while trying to delete repo :(", getDefaultSendOptions())
	}

	repo, err := botService.DatabaseService.GetRepoWhereId(id)
	if err != nil {
		botService.Logger.Error("An error occured while trying to GET repo", zap.Error(err))
		return err
	}

	err = botService.DatabaseService.DeleteRepoWhereId(id)
	if err != nil {
		botService.Logger.Error("An error occured while trying to DeleteRepoWhereId", zap.Error(err))
		return err
	}

	err = state.Finish(c.Data() != "")
	if err != nil {
		botService.Logger.Error("An error occured while trying to finish RemoveRepoState ", zap.Error(err))
		return err
	}

	_, err = botService.bot.Edit(c.Callback(), &telebot.ReplyMarkup{})
	if err != nil {
		botService.Logger.Error("An error occured while trying to remove inline keyboard markup", zap.Error(err))
		return err
	}

	message_text := fmt.Sprintf("Repo <code>%s:%s/%s</code> successfully removed", repo.Host, repo.Owner, repo.Repo)

	sendOptions := getDefaultSendOptions()
	sendOptions.ReplyMarkup = getStartKeyboard()

	return c.Send(message_text, sendOptions)
}

func (botService *BotService) AddRepoHandler(c telebot.Context, state fsm.Context) error {
	botService.Logger.Info(
		"Received /add_repo, handling command...",
		zap.String("Sender username", c.Sender().Username),
	)

	err := state.Set(AddRepoState)
	if err != nil {
		botService.Logger.Error(
			fmt.Sprintf("An error occured while setting a state %s", AddRepoState),
			zap.Error(err),
			zap.String("Sender username", c.Sender().Username),
		)
	}

	err = c.Send("No problem! Just send me a link to the repository (GitHub only for now)", getHomeKeyboard())

	return err
}

func (botService *BotService) OnRepoEntered(c telebot.Context, state fsm.Context) error {
	var githubErrResponse *github.ErrorResponse

	repo, err := getRepoFromMessage(c.Message(), botService.Fetcher)

	if err != nil {
		messageText := ""
		if errors.Is(err, ErrHostIsIncorrect) {
			messageText = fmt.Sprintf("Host <code>%s</code> is incorrect or doesn't support yet :(", repo.Host)
		} else if errors.Is(err, ErrCantParseRepo) {
			messageText = "Sorry, but I can't add the repository :(\n" +
				"Please enter the repo in the format 'https://<u>host</u>/<u>owner</u>/<u>repo</u>'"
		} else if errors.Is(err, fetcher.ErrNoTagsInRepo) {
			messageText = "There is no releases/tags in the repository:(\n" +
				"Maybe add another repo?"
		} else if errors.As(err, &githubErrResponse) {
			messageText = fmt.Sprintf("Repo <code>%s</code> is incorrect/private or doesn't exist :(", repo.Repo)
		} else {
			botService.Logger.Fatal(
				"An error occured while validating repo",
				zap.Error(err),
			)
		}
		err := c.Send(messageText, &telebot.SendOptions{ParseMode: telebot.ModeHTML})

		return err
	}

	isRepoExist, err := botService.DatabaseService.IsRepoAlreadyExist(repo)
	if err != nil {
		botService.Logger.Fatal(
			"An error occured while trying to db.Get repo",
			zap.Error(err),
			zap.String("Repo owner", repo.Owner),
			zap.String("Repo", repo.Repo),
		)
	}

	if isRepoExist {
		messageText := fmt.Sprintf("<code>%s:%s/%s</code> is already exist", repo.Host, repo.Owner, repo.Repo)
		err := c.Send(messageText, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
		if err != nil {
			botService.Logger.Error(
				"An error occured while trying to send 'repo is already exist' for user",
				zap.Error(err),
			)
		}
		return err
	}

	err = botService.DatabaseService.AddRepo(repo)
	if err != nil {
		botService.Logger.Error(
			"An error occured while trying to add a new repo",
			zap.Error(err),
		)
		return err
	}

	err = state.Finish(c.Data() != "")
	if err != nil {
		botService.Logger.Error(
			"An error occured while trying to finish AddRepo state",
			zap.Error(err),
		)
		return err
	}

	messageText := fmt.Sprintf("Repo <code>%s:%s/%s</code> successfully added!", repo.Host, repo.Owner, repo.Repo)

	sendOptions := getDefaultSendOptions()
	sendOptions.ReplyMarkup = getStartKeyboard()

	err = c.Send(messageText, sendOptions)
	if err != nil {
		botService.Logger.Error(
			"An error occured while sending message to user",
			zap.Error(err),
		)
	}

	return err
}
