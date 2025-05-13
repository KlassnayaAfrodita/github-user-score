package services

import (
	"context"
	"fmt"
)

func (service *ScoringManagerService) MarkExpiredApplications(ctx context.Context, maxAgeMinutes int) error {
	expiredIDs, err := service.repo.GetExpiredApplications(ctx, maxAgeMinutes)
	if err != nil {
		return fmt.Errorf("ScoringManagerService.MarkExpiredApplications: %w", err)
	}

	if len(expiredIDs) == 0 {
		return nil
	}

	err = service.repo.MarkExpiredApplications(ctx, expiredIDs)
	if err != nil {
		return fmt.Errorf("ScoringManagerService.MarkExpiredApplications: %w", err)
	}

	return nil
}
