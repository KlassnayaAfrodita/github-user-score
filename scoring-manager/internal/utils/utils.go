package utils

import (
	collector "github.com/KlassnayaAfrodita/github-user-score/collector/pkg/pb/collector/api"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
)

func ToUserStats(response *collector.CollectUserStatsResponse) repository.UserStats {
	return repository.UserStats{
		UserID:  response.UserID,
		Repos:   response.Repos,
		Stars:   response.Stars,
		Forks:   response.Forks,
		Commits: response.Commits,
	}
}

func IsEmptyScoringApplication(s repository.ScoringApplication) bool {
	return s.ApplicationID == 0 && s.UserID == 0
}
