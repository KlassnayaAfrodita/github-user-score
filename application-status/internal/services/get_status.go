package services

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/application-status/internal/client/grpc"
)

func (service *ApplicationStatusClient) GetStatus(ctx context.Context, applicationID int64) (grpc.GetScoreResponse, error) {
	return service.scoringManager.GetStatus(ctx, applicationID)
}
