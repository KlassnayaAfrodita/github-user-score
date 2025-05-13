package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/mocks"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/services"
)

func TestMarkExpiredApplications_NoExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	maxAge := 60

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	mockRepo.
		EXPECT().
		GetExpiredApplications(ctx, maxAge).
		Return([]int64{}, nil)

	service := services.NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	err := service.MarkExpiredApplications(ctx, maxAge)
	require.NoError(t, err)
}

func TestMarkExpiredApplications_WithExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	maxAge := 60
	expired := []int64{1, 2, 3}

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	gomock.InOrder(
		mockRepo.EXPECT().
			GetExpiredApplications(ctx, maxAge).
			Return(expired, nil),
		mockRepo.EXPECT().
			MarkExpiredApplications(ctx, expired).
			Return(nil),
	)

	service := services.NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	err := service.MarkExpiredApplications(ctx, maxAge)
	require.NoError(t, err)
}

func TestMarkExpiredApplications_GetExpiredError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	maxAge := 60
	expectedErr := errors.New("db failed")

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	mockRepo.
		EXPECT().
		GetExpiredApplications(ctx, maxAge).
		Return(nil, expectedErr)

	service := services.NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	err := service.MarkExpiredApplications(ctx, maxAge)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db failed")
}

func TestMarkExpiredApplications_MarkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	maxAge := 60
	expired := []int64{100, 200}
	expectedErr := errors.New("update failed")

	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	gomock.InOrder(
		mockRepo.EXPECT().
			GetExpiredApplications(ctx, maxAge).
			Return(expired, nil),
		mockRepo.EXPECT().
			MarkExpiredApplications(ctx, expired).
			Return(expectedErr),
	)

	service := services.NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	err := service.MarkExpiredApplications(ctx, maxAge)
	require.Error(t, err)
	require.Contains(t, err.Error(), "update failed")
}
