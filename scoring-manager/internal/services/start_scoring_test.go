package services

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/mocks"

	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestScoringManagerService_StartScoring_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	username := "testuser"
	expectedUserID := int64(123)
	expectedAppID := int64(1)

	userStats := repository.UserStats{
		UserID:  int32(expectedUserID),
		Repos:   5,
		Stars:   10,
		Forks:   3,
		Commits: 50,
	}

	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockCollector.EXPECT().
		GetUserStats(ctx, username).
		Return(userStats, nil)

	expectedApp := repository.ScoringApplication{
		UserID: int32(expectedUserID),
		Status: repository.StatusInitial,
		Score:  nil,
	}
	createdApp := expectedApp
	createdApp.ApplicationID = expectedAppID

	mockRepo.EXPECT().
		CreateScoringApplication(ctx, expectedApp).
		Return(createdApp, nil)

	expectedMsg := kafka.ScoringRequestMessage{
		ApplicationID: strconv.Itoa(int(expectedAppID)),
		UserID:        int(expectedUserID),
		Repos:         int(userStats.Repos),
		Stars:         int(userStats.Stars),
		Forks:         int(userStats.Forks),
		Commits:       int(userStats.Commits),
	}
	mockKafka.EXPECT().
		PublishScoringRequest(ctx, expectedMsg).
		Return(nil)

	appID, err := service.StartScoring(ctx, username)

	require.NoError(t, err)
	require.Equal(t, strconv.Itoa(int(expectedAppID)), appID)
}

func TestScoringManagerService_StartScoring_GetUserStatsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	username := "testuser"
	expectedErr := errors.New("user not found")

	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockCollector.EXPECT().
		GetUserStats(ctx, username).
		Return(repository.UserStats{}, expectedErr)

	appID, err := service.StartScoring(ctx, username)

	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErr.Error())
	require.Empty(t, appID)
}

func TestScoringManagerService_StartScoring_CreateApplicationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	username := "testuser"
	userStats := repository.UserStats{UserID: 123}
	expectedErr := errors.New("database error")

	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockCollector.EXPECT().
		GetUserStats(ctx, username).
		Return(userStats, nil)

	expectedApp := repository.ScoringApplication{
		UserID: userStats.UserID,
		Status: repository.StatusInitial,
		Score:  nil,
	}

	mockRepo.EXPECT().
		CreateScoringApplication(ctx, expectedApp).
		Return(repository.ScoringApplication{}, expectedErr)

	appID, err := service.StartScoring(ctx, username)

	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErr.Error())
	require.Empty(t, appID)
}

func TestScoringManagerService_StartScoring_PublishError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	username := "testuser"
	expectedAppID := int64(1)
	expectedErr := errors.New("kafka error")

	userStats := repository.UserStats{
		UserID:  123,
		Repos:   2,
		Stars:   5,
		Forks:   1,
		Commits: 10,
	}

	mockCollector := mocks.NewMockCollectorClientInterface(ctrl)
	mockRepo := mocks.NewMockScoringRepositoryInterface(ctrl)
	mockKafka := mocks.NewMockScoringKafkaClient(ctrl)

	service := NewScoringManagerService(mockRepo, mockCollector, mockKafka)

	mockCollector.EXPECT().
		GetUserStats(ctx, username).
		Return(userStats, nil)

	expectedApp := repository.ScoringApplication{
		UserID: userStats.UserID,
		Status: repository.StatusInitial,
		Score:  nil,
	}
	createdApp := expectedApp
	createdApp.ApplicationID = expectedAppID

	mockRepo.EXPECT().
		CreateScoringApplication(ctx, expectedApp).
		Return(createdApp, nil)

	expectedMsg := kafka.ScoringRequestMessage{
		ApplicationID: strconv.Itoa(int(expectedAppID)),
		UserID:        int(userStats.UserID),
		Repos:         int(userStats.Repos),
		Stars:         int(userStats.Stars),
		Forks:         int(userStats.Forks),
		Commits:       int(userStats.Commits),
	}

	mockKafka.EXPECT().
		PublishScoringRequest(ctx, expectedMsg).
		Return(expectedErr)

	appID, err := service.StartScoring(ctx, username)

	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErr.Error())
	require.Equal(t, strconv.Itoa(int(expectedAppID)), appID)
}
