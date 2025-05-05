package controllers

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/services/utils"
	pb "github.com/KlassnayaAfrodita/github-user-score/collector/pb"
)

func (h *CollectorHandler) CollectUserStats(ctx context.Context, req *pb.CollectUserStatsRequest) (*pb.CollectUserStatsResponse, error) {
	stats, err := h.service.CollectStats(ctx, req.GetUsername())
	if err != nil {
		return &pb.CollectUserStatsResponse{}, nil
	}

	return utils.ToProtoStats(stats), nil
}
