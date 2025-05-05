package services

import (
	"context"
	github "github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/github-api"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
)

type CollectorServiceInterface interface {
	CollectStats(ctx context.Context, username string) (repository.Stats, error)
}

type CollectorService struct {
	repo         repository.CollectorRepositoryInterface
	githubClient github.GitHubClientInterface
}

func NewCollectorService(repo repository.CollectorRepositoryInterface, client github.GitHubClientInterface) *CollectorService {
	return &CollectorService{
		repo:         repo,
		githubClient: client,
	}
}
