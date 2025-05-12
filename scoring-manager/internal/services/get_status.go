package services

import (
	"context"
	"fmt"
)

func (service *ScoringManagerService) GetStatus(ctx context.Context, applicationID int) (ScoringStatus, error) {
	scoringApplication, err := service.repo.GetScoringApplicationByID(ctx, applicationID)
	if err != nil {
		// возвращать более информотивную ошибку. типа напоминание пользователю, что нужно еще проскорить, потому что
		// либо неправильный applicationID, либо заявка зависла и мы ее удалили
		return ScoringStatus{}, fmt.Errorf("ScoringManagerService.GetStatus: %w", err)
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
