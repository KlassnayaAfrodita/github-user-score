package services

import (
	"context"
	"errors"
	"testing"

	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
	"github.com/KlassnayaAfrodita/github-user-score/collector/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCollectorService_CollectStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCollectorRepositoryInterface(ctrl)
	mockGitHub := mocks.NewMockGitHubClientInterface(ctrl)

	service := NewCollectorService(mockRepo, mockGitHub)
	ctx := context.Background()
	username := "testuser"

	t.Run("success - user exists", func(t *testing.T) {
		user := &repository.User{ID: 1, Username: username}
		stats := repository.Stats{Commits: 10}

		mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(user, nil)
		mockGitHub.EXPECT().GetStats(ctx, username).Return(stats, nil)
		mockRepo.EXPECT().SaveUserStats(ctx, repository.Stats{UserID: 1, Commits: 10}).Return(nil)

		err := service.CollectStats(ctx, username)
		require.NoError(t, err)
	})

	t.Run("success - user does not exist", func(t *testing.T) {
		user := &repository.User{ID: 2, Username: username}
		stats := repository.Stats{Commits: 5}

		mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(ctx, username).Return(user, nil)
		mockGitHub.EXPECT().GetStats(ctx, username).Return(stats, nil)
		mockRepo.EXPECT().SaveUserStats(ctx, repository.Stats{UserID: 2, Commits: 5}).Return(nil)

		err := service.CollectStats(ctx, username)
		require.NoError(t, err)
	})

	t.Run("error getting user", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(nil, errors.New("db error"))

		err := service.CollectStats(ctx, username)
		require.ErrorContains(t, err, "service.CollectStats")
	})

	t.Run("error creating user", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(ctx, username).Return(nil, errors.New("insert error"))

		err := service.CollectStats(ctx, username)
		require.ErrorContains(t, err, "service.CollectStats")
	})

	t.Run("error getting stats", func(t *testing.T) {
		user := &repository.User{ID: 3, Username: username}
		mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(user, nil)
		mockGitHub.EXPECT().GetStats(ctx, username).Return(repository.Stats{}, errors.New("api error"))

		err := service.CollectStats(ctx, username)
		require.ErrorContains(t, err, "service.CollectStats")
	})

	t.Run("error saving stats", func(t *testing.T) {
		user := &repository.User{ID: 4, Username: username}
		stats := repository.Stats{Commits: 8}

		mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(user, nil)
		mockGitHub.EXPECT().GetStats(ctx, username).Return(stats, nil)
		mockRepo.EXPECT().SaveUserStats(ctx, repository.Stats{UserID: 4, Commits: 8}).Return(errors.New("save error"))

		err := service.CollectStats(ctx, username)
		require.ErrorContains(t, err, "service.CollectStats")
	})
}
