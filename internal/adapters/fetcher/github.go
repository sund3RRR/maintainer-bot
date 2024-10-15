package fetcher

import (
	"context"
	"fmt"

	"github.com/sund3RRR/maintainer-bot/internal/adapters/db"

	"github.com/google/go-github/v57/github"
	"go.uber.org/zap"
)

type GithubFetcher struct {
	Client *github.Client
}

func (f *GithubFetcher) FetchRepo(repo *db.Repo, logger *zap.Logger) *RepoMessage {
	newTagName, body, link := "", "", ""
	if repo.IsRelease {
		release, _, err := f.Client.Repositories.GetLatestRelease(context.Background(), repo.Owner, repo.Repo)
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
		tags, _, err := f.Client.Repositories.ListTags(context.Background(), repo.Owner, repo.Repo, &github.ListOptions{Page: 0})
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

func (f *GithubFetcher) GetGithubRepository(repo *db.Repo) (*github.Repository, error) {
	repository, _, err := f.Client.Repositories.Get(context.Background(), repo.Owner, repo.Repo)

	return repository, err
}

func (f *GithubFetcher) GetLatestReleaseTagName(repo *db.Repo) (string, error) {
	lastRelease, _, err := f.Client.Repositories.GetLatestRelease(context.Background(), repo.Owner, repo.Repo)
	if err != nil {
		return "", err
	}
	tagName := lastRelease.GetTagName()

	return tagName, nil
}

func (f *GithubFetcher) GetLatestTagName(repo *db.Repo) (string, error) {
	tags, _, err := f.Client.Repositories.ListTags(context.Background(), repo.Owner, repo.Repo, &github.ListOptions{Page: 1})
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "", ErrNoTagsInRepo
	}

	tagName := tags[0].GetName()

	return tagName, nil
}
