package services

import (
	"context"
	"errors"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestConsumingScoringResults_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	msg := kafka.ScoringResultMessage{ApplicationID: 1, UserID: 2, Score: 100}

	mockKafka.EXPECT().ConsumeScoringResults(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, handler func(kafka.ScoringResultMessage) error) error {
			return handler(msg)
		})

	mockRepo.EXPECT().SaveScoringApplicationResult(ctx, repository.ScoringApplication{
		ApplicationID: msg.ApplicationID,
		UserID:        int32(msg.UserID),
		Score:         &msg.Score,
	}).Return(nil)

	mockRepo.EXPECT().UpdateScoringApplicationStatus(ctx, msg.ApplicationID, repository.StatusSuccess).
		Return(nil)

	err := service.ConsumingScoringResults(ctx)
	assert.NoError(t, err)
}

func TestConsumingScoringResults_SaveResultError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	msg := kafka.ScoringResultMessage{ApplicationID: 1, UserID: 2, Score: 100}

	mockKafka.EXPECT().ConsumeScoringResults(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, handler func(kafka.ScoringResultMessage) error) error {
			return handler(msg)
		})

	mockRepo.EXPECT().SaveScoringApplicationResult(ctx, gomock.Any()).
		Return(errors.New("failed to save"))

	err := service.ConsumingScoringResults(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save")
}

func TestConsumingScoringResults_UpdateStatusError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	msg := kafka.ScoringResultMessage{ApplicationID: 1, UserID: 2, Score: 100}

	mockKafka.EXPECT().ConsumeScoringResults(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, handler func(kafka.ScoringResultMessage) error) error {
			mockRepo.EXPECT().SaveScoringApplicationResult(ctx, gomock.Any()).Return(nil)
			mockRepo.EXPECT().UpdateScoringApplicationStatus(ctx, msg.ApplicationID, repository.StatusSuccess).
				Return(errors.New("update error"))
			return handler(msg)
		})

	err := service.ConsumingScoringResults(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
}

func TestConsumingScoringResults_ConsumeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockKafka.EXPECT().ConsumeScoringResults(gomock.Any(), gomock.Any()).
		Return(errors.New("kafka consumer failure"))

	err := service.ConsumingScoringResults(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kafka consumer failure")
}
