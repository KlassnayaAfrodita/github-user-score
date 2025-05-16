package grpc

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/pkg/pb/scoring-manager/api"
	"google.golang.org/grpc"
)

type GetScoreResponse struct {
	Status  ScoringStatus
	Scoring int32
}

type ScoringManagerClientInterface interface {
	GetStatus(ctx context.Context, applicationID int64) (GetScoreResponse, error)
}

type ScoringManagerClient struct {
	conn    *grpc.ClientConn
	service api.ScoringManagerServiceClient
}
