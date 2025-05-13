package services

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
)

func (service *ScoringManagerService) ConsumingScoringResults(ctx context.Context) error {
	handler := func(msg kafka.ScoringResultMessage) error {
		app := repository.ScoringApplication{
			ApplicationID: msg.ApplicationID,
			UserID:        int32(msg.UserID),
			Score:         &msg.Score,
		}

		if err := service.repo.SaveScoringApplicationResult(ctx, app); err != nil {
			return fmt.Errorf("StartConsumingScoringResults: failed to save result: %w", err)
		}

		if err := service.repo.UpdateScoringApplicationStatus(ctx, msg.ApplicationID, repository.StatusSuccess); err != nil {
			return fmt.Errorf("StartConsumingScoringResults: failed to update status: %w", err)
		}

		return nil
	}

	return service.kafka.ConsumeScoringResults(ctx, handler)
}
