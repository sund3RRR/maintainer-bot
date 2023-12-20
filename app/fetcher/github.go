package fetcher

import (
	repodb "app/db"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v57/github"
	"go.uber.org/zap"
)

func formatTitle(repo *repodb.Repo, newTagName string) string {
	return fmt.Sprintf(
		"<b>%s</b> / <b>%s</b> <code>%s</code> -> <code>%s</code>",
		repo.Owner,
		repo.Repo,
		repo.LastTag,
		newTagName,
	)
}
func formatReleaseBody(text string) string {
	splitted := strings.Split(text, "\n")
	var result string

	for _, line := range splitted {
		is_heading := false
		for len(line) > 0 && line[0] == '#' {
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimSpace(line)
			is_heading = true
		}
		if is_heading {
			result += fmt.Sprintf("<b><u>%s</u></b>\n", line)
		} else {
			result += line + "\n"
		}

	}

	result = replaceMdToHtml(result, "`", "code")
	result = replaceMdToHtml(result, `\*\*`, "b")
	result = replaceMdToHtml(result, `\*`, "i")

	return result
}

func replaceMdToHtml(text, mdSymbol, htmlTag string) string {
	// Define a regular expression to find words enclosed in backticks
	pattern := fmt.Sprintf(`%s([^`+"`"+`]+)%s`, mdSymbol, mdSymbol)

	// Replace all occurrences of the pattern with <code>word</code>
	re := regexp.MustCompile(pattern)
	result := re.ReplaceAllString(text, fmt.Sprintf(`<%s>$1</%s>`, htmlTag, htmlTag))

	return result
}

func FetchGithubRepo(repo *repodb.Repo, client *github.Client, logger *zap.Logger) *RepoMessage {
	newTagName, body, link := "", "", ""
	if repo.IsRelease {
		release, _, err := client.Repositories.GetLatestRelease(context.Background(), repo.Owner, repo.Repo)
		if err != nil {
			logger.Error(
				"An error occured while getting repositiry",
				zap.Error(err),
				zap.String("Repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Repo)),
			)
		}
		newTagName = release.GetTagName()
		body = formatReleaseBody(*release.Body)
		link = fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", repo.Owner, repo.Repo, newTagName)
	} else {
		tags, _, err := client.Repositories.ListTags(context.Background(), repo.Owner, repo.Repo, &github.ListOptions{Page: 0})
		if err != nil {
			logger.Error(
				"An error occured while getting tags",
				zap.Error(err),
				zap.String("Repo", fmt.Sprintf("%s/%s", repo.Owner, repo.Repo)),
			)
		}
		newTagName = tags[0].GetName()
		body = tags[0].GetCommit().GetMessage()
		link = fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", repo.Owner, repo.Repo, newTagName)
	}

	title := formatTitle(repo, newTagName)

	return &RepoMessage{
		ChatID: repo.ChatID,
		Title:  title,
		Text:   body,
		Link:   link,
		NewTag: newTagName,
	}
}
