package services

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	"strconv"
)

func (service *ScoringManagerService) StartScoring(ctx context.Context, username string) (string, error) {
	userStats, err := service.collector.GetUserStats(ctx, username)
	if err != nil {
		return "", fmt.Errorf("ScoringManagerService.StartScoring: %w", err)
	}

	var zeroScore *int

	scoringApplication := repository.ScoringApplication{
		UserID: userStats.UserID,
		Status: repository.StatusInitial,
		Score:  zeroScore,
	}

	scoringApplication, err = service.repo.CreateScoringApplication(ctx, scoringApplication)
	if err != nil {
		return "", fmt.Errorf("ScoringManagerService.StartScoring: %w", err)
	}

	msg := kafka.ScoringRequestMessage{
		ApplicationID: strconv.Itoa(int(scoringApplication.ApplicationID)),
		UserID:        int(scoringApplication.UserID),
		Repos:         int(userStats.Repos),
		Stars:         int(userStats.Stars),
		Forks:         int(userStats.Forks),
		Commits:       int(userStats.Commits),
	}

	err = service.kafka.PublishScoringRequest(ctx, msg)
	if err != nil {
		return msg.ApplicationID, fmt.Errorf("ScoringManagerService.StartScoring: %w", err)
	}

	return msg.ApplicationID, nil
}
