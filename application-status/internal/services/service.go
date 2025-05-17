package services

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/application-status/internal/client/grpc"
)

type ApplicationStatusClientInterface interface {
	GetStatus(ctx context.Context, applicationID int64) (grpc.GetScoreResponse, error)
}

type ApplicationStatusClient struct {
	scoringManager grpc.ScoringManagerClientInterface
}

func NewApplicationStatusClient(scoringManager grpc.ScoringManagerClientInterface) ApplicationStatusClientInterface {
	return ApplicationStatusClient{scoringManager: scoringManager}
}
