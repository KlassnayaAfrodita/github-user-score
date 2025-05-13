package services

import "fmt"

func (service *ScoringManagerService) MarkExpiredApplications(ctx context.Context, maxAgeMinutes int) (int64, error) {
	affected, err := service.repo.MarkExpiredApplications(ctx, maxAgeMinutes)
	if err != nil {
		return 0, fmt.Errorf("ScoringManagerService.CleanupStuckApplications: %w", err)
	}
	return affected, nil
}
