package services

import (
	"context"
	"fmt"
	github "github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/github-api"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
)

type CollectorServiceInterface interface {
	CollectStats(ctx context.Context, username string) error
}

type CollectorService struct {
	repo   repository.CollectorRepositoryInterface
	client github.GitHubClientInterface
}

func NewCollectorService(repo repository.CollectorRepositoryInterface, client github.GitHubClientInterface) *CollectorService {
	return &CollectorService{
		repo:   repo,
		client: client,
	}
}

func (s *CollectorService) CollectStats(ctx context.Context, username string) (repository.Stats, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return repository.Stats{}, fmt.Errorf("service.CollectStats: %w", err)
	}
	if user == nil {
		user, err = s.repo.CreateUser(ctx, username)
		if err != nil {
			return repository.Stats{}, fmt.Errorf("service.CollectStats: %w", err)
		}
	}

	stats, err := s.client.GetStats(ctx, username)
	if err != nil {
		return repository.Stats{}, fmt.Errorf("service.CollectStats: %w", err)
	}

	stats.UserID = user.ID

	if err := s.repo.SaveUserStats(ctx, stats); err != nil {
		return repository.Stats{}, fmt.Errorf("service.CollectStats: %w", err)
	}

	return stats, nil
}
