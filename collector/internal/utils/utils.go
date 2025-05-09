package utils

import (
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
	pb "github.com/KlassnayaAfrodita/github-user-score/collector/pkg/pb/collector/api"
)

func ToProtoStats(stats repository.Stats) *pb.CollectUserStatsResponse {
	return &pb.CollectUserStatsResponse{
		UserID:  int32(stats.UserID),
		Repos:   int32(stats.Repos),
		Stars:   int32(stats.Stars),
		Forks:   int32(stats.Forks),
		Commits: int32(stats.Commits),
	}
}
