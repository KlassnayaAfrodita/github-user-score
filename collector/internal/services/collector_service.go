package services

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
)

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

	stats, err := s.githubClient.GetStats(ctx, username)
	if err != nil {
		return repository.Stats{}, fmt.Errorf("service.CollectStats: %w", err)
	}

	stats.UserID = user.ID

	if err := s.repo.SaveUserStats(ctx, stats); err != nil {
		return repository.Stats{}, fmt.Errorf("service.CollectStats: %w", err)
	}

	return stats, nil
}
