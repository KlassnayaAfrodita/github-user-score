package controllers

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/utils"
	pb "github.com/KlassnayaAfrodita/github-user-score/collector/pkg/pb/collector/api"
)

func (h *CollectorHandler) CollectUserStats(ctx context.Context, req *pb.CollectUserStatsRequest) (*pb.CollectUserStatsResponse, error) {
	stats, err := h.service.CollectStats(ctx, req.GetUsername())
	if err != nil {
		return &pb.CollectUserStatsResponse{}, err
	}

	return utils.ToProtoStats(stats), nil
}
