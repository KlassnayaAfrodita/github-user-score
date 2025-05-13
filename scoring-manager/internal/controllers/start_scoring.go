package controllers

import (
	"context"

	pb "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/pkg/pb/scoring-manager/api"
)

func (c *ScoringManagerHandler) StartScoring(ctx context.Context, req *pb.StartScoringRequest) (*pb.StartScoringResponse, error) {
	appID, err := c.service.StartScoring(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.StartScoringResponse{
		ApplicationId: appID,
	}, nil
}
