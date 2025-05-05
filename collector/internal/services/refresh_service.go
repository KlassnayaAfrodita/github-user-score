package services

import (
	"context"
	"time"
)

func (s *CollectorService) RefreshOutdatedStats(ctx context.Context, threshold time.Duration) error {
	users, err := s.repo.GetOutdatedUsers(ctx, threshold)
	if err != nil {
		return err
	}

	for _, user := range users {
		stats, err := s.githubClient.GetStats(ctx, user.Username)
		if err != nil {
			return err
		}
		stats.UserID = user.ID
		if err := s.repo.SaveUserStats(ctx, stats); err != nil {
			return err
		}
	}

	return nil
}
