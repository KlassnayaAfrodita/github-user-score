package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestScoringManagerService_GetStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	appID := 42
	score := int(85)

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockRepo.EXPECT().
		GetScoringApplicationByID(ctx, appID).
		Return(repository.ScoringApplication{
			ApplicationID: int64(appID),
			Status:        1,
			Score:         &score,
		}, nil)

	status, err := service.GetStatus(ctx, appID)

	require.NoError(t, err)
	require.Equal(t, int32(1), status.status)
	require.Equal(t, int32(score), status.scoringResult)
}

func TestScoringManagerService_GetStatus_NotFinished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	appID := 42

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockRepo.EXPECT().
		GetScoringApplicationByID(ctx, appID).
		Return(repository.ScoringApplication{
			ApplicationID: int64(appID),
			Status:        0,
			Score:         nil,
		}, nil)

	status, err := service.GetStatus(ctx, appID)

	require.Error(t, err)
	require.Contains(t, err.Error(), "scoring is not finished yet")
	require.Equal(t, ScoringStatus{}, status)
}

func TestScoringManagerService_GetStatus_ErrorStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	appID := 42

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockRepo.EXPECT().
		GetScoringApplicationByID(ctx, appID).
		Return(repository.ScoringApplication{
			ApplicationID: int64(appID),
			Status:        2,
			Score:         nil,
		}, nil)

	status, err := service.GetStatus(ctx, appID)

	require.Error(t, err)
	require.Contains(t, err.Error(), "scoring ended in error")
	require.Equal(t, ScoringStatus{}, status)
}

func TestScoringManagerService_GetStatus_UnknownStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	appID := 42

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockRepo.EXPECT().
		GetScoringApplicationByID(ctx, appID).
		Return(repository.ScoringApplication{
			ApplicationID: int64(appID),
			Status:        999,
			Score:         nil,
		}, nil)

	status, err := service.GetStatus(ctx, appID)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown error")
	require.Equal(t, ScoringStatus{}, status)
}

func TestScoringManagerService_GetStatus_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	appID := 42
	expectedErr := errors.New("db connection failed")

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockRepo.EXPECT().
		GetScoringApplicationByID(ctx, appID).
		Return(repository.ScoringApplication{}, expectedErr)

	status, err := service.GetStatus(ctx, appID)

	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("ScoringManagerService.GetStatus: %s", expectedErr.Error()))
	require.Equal(t, ScoringStatus{}, status)
}
