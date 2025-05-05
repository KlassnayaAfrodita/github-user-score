package controllers

import (
	"context"
	pb "github.com/KlassnayaAfrodita/github-user-score/collector/internal/pb/collector/api"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/utils"
)

func (h *CollectorHandler) CollectUserStats(ctx context.Context, req *pb.CollectUserStatsRequest) (*pb.CollectUserStatsResponse, error) {
	stats, err := h.service.CollectStats(ctx, req.GetUsername())
	if err != nil {
		return &pb.CollectUserStatsResponse{}, nil
	}

	return utils.ToProtoStats(stats), nil
}
