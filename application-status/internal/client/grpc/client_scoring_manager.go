package grpc

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/pkg/pb/scoring-manager/api"
	"google.golang.org/grpc"
	"time"
)

type ScoringStatus int

const (
	//INITIAL
	StatusInitial ScoringStatus = 0
	//SUCCESS
	StatusSuccess ScoringStatus = 1
	//FAILED
	StatusFailed ScoringStatus = 2
)

func NewScoringManagerClient(address string, timeout time.Duration) (ScoringManagerClientInterface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &ScoringManagerClient{
		conn:    conn,
		service: api.NewScoringManagerServiceClient(conn),
	}, nil
}

func (client *ScoringManagerClient) GetStatus(ctx context.Context, applicationID int64) (GetScoreResponse, error) {
	req := &api.GetStatusRequest{ApplicationId: applicationID}
	userScore, err := client.service.GetStatus(ctx, req)
	if err != nil {
		return GetScoreResponse{}, fmt.Errorf("ScoringMangerClient.GetStatus: %w", err)
	}

	return ToUserScore(userScore), nil
}

func ToUserScore(response *api.GetStatusResponse) GetScoreResponse {
	return GetScoreResponse{
		Status:  ScoringStatus(response.Status),
		Scoring: response.Scoring,
	}
}
