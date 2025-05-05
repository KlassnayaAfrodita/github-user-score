package github_api

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
)

type GitHubClientInterface interface {
	GetStats(ctx context.Context, username string) (repository.Stats, error)
}

type GitHubClient struct{}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{}
}

// пока заглушка
func (g *GitHubClient) GetStats(ctx context.Context, username string) (repository.Stats, error) {
	_ = username
	_ = ctx
	return repository.Stats{
		UserID:  1,
		Repos:   1,
		Stars:   1,
		Forks:   1,
		Commits: 1,
	}, nil
}
