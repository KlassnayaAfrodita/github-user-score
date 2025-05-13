package services

import (
	"context"
	collectorClient "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/grpc"
	kafkaClient "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
)

type ScoringStatus struct {
	Status        int32
	ScoringResult int32
}

type ScoringManagerServiceInterface interface {
	StartScoring(ctx context.Context, username string) (int64, error)
	GetStatus(ctx context.Context, applicationID int) (ScoringStatus, error)
	MarkExpiredApplications(ctx context.Context, maxAgeMinutes int) error
	ConsumingScoringResults(ctx context.Context) error
}

type ScoringManagerService struct {
	repo      repository.ScoringRepositoryInterface
	collector collectorClient.CollectorClientInterface
	kafka     kafkaClient.ScoringKafkaClient
}

func NewScoringManagerService(repo repository.ScoringRepositoryInterface,
	collector collectorClient.CollectorClientInterface,
	kafka kafkaClient.ScoringKafkaClient) ScoringManagerServiceInterface {
	return &ScoringManagerService{repo: repo, collector: collector, kafka: kafka}
}
