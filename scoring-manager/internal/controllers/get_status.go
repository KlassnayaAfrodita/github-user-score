package controllers

import (
	"context"
	pb "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/pkg/pb/scoring-manager/api"
)

func (c *ScoringManagerHandler) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	scoringStatus, err := c.service.GetStatus(ctx, int(req.ApplicationId))
	if err != nil {
		return nil, err
	}

	resp := &pb.GetStatusResponse{
		Status: pb.ScoringStatus(scoringStatus.Status),
	}

	if scoringStatus.ScoringResult != 0 {
		resp.Scoring = scoringStatus.ScoringResult
	}

	return resp, nil
}
