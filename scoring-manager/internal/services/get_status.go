package services

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/utils"
)

func (service *ScoringManagerService) GetStatus(ctx context.Context, applicationID int) (ScoringStatus, error) {
	scoringApplication, err := service.repo.GetScoringApplicationByID(ctx, applicationID)
	if err != nil {
		return ScoringStatus{}, fmt.Errorf("ScoringManagerService.GetStatus: %w", err)
	}

	if utils.IsEmptyScoringApplication(scoringApplication) {
		return ScoringStatus{}, fmt.Errorf("Заявки нет в скоринге")
	}

	switch scoringApplication.Status {
	case 0:
		return ScoringStatus{}, fmt.Errorf("scoring is not finished yet")
	case 1:
		return ScoringStatus{
			status:        int32(scoringApplication.Status),
			scoringResult: int32(*scoringApplication.Score),
		}, nil
	case 2:
		return ScoringStatus{}, fmt.Errorf("scoring ended in error")
	default:
		return ScoringStatus{}, fmt.Errorf("unknown error")
	}
}
