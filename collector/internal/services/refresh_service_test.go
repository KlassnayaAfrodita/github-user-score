package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
	"github.com/KlassnayaAfrodita/github-user-score/collector/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCollectorService_RefreshOutdatedStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCollectorRepositoryInterface(ctrl)
	mockGitHub := mocks.NewMockGitHubClientInterface(ctrl)

	service := NewCollectorService(mockRepo, mockGitHub)
	ctx := context.Background()
	threshold := time.Hour

	t.Run("success - users with outdated stats", func(t *testing.T) {
		users := []repository.User{
			{ID: 1, Username: "user1"},
			{ID: 2, Username: "user2"},
		}
		stats1 := repository.Stats{Commits: 10}
		stats2 := repository.Stats{Commits: 20}

		mockRepo.EXPECT().GetOutdatedUsers(ctx, threshold).Return(users, nil)
		mockGitHub.EXPECT().GetStats(ctx, "user1").Return(stats1, nil)
		mockGitHub.EXPECT().GetStats(ctx, "user2").Return(stats2, nil)
		mockRepo.EXPECT().SaveUserStats(ctx, repository.Stats{UserID: 1, Commits: 10}).Return(nil)
		mockRepo.EXPECT().SaveUserStats(ctx, repository.Stats{UserID: 2, Commits: 20}).Return(nil)

		err := service.RefreshOutdatedStats(ctx, threshold)
		require.NoError(t, err)
	})

	t.Run("success - no outdated users", func(t *testing.T) {
		mockRepo.EXPECT().GetOutdatedUsers(ctx, threshold).Return([]repository.User{}, nil)

		err := service.RefreshOutdatedStats(ctx, threshold)
		require.NoError(t, err)
	})

	t.Run("error getting outdated users", func(t *testing.T) {
		mockRepo.EXPECT().GetOutdatedUsers(ctx, threshold).Return(nil, errors.New("db error"))

		err := service.RefreshOutdatedStats(ctx, threshold)
		require.ErrorContains(t, err, "db error")
	})

	t.Run("error getting stats from GitHub", func(t *testing.T) {
		users := []repository.User{
			{ID: 1, Username: "user1"},
			{ID: 2, Username: "user2"},
		}
		mockRepo.EXPECT().GetOutdatedUsers(ctx, threshold).Return(users, nil)
		mockGitHub.EXPECT().GetStats(ctx, "user1").Return(repository.Stats{}, errors.New("api error"))

		err := service.RefreshOutdatedStats(ctx, threshold)
		require.ErrorContains(t, err, "api error")
	})

	t.Run("error saving stats", func(t *testing.T) {
		users := []repository.User{
			{ID: 1, Username: "user1"},
			{ID: 2, Username: "user2"},
		}
		stats1 := repository.Stats{Commits: 10}

		mockRepo.EXPECT().GetOutdatedUsers(ctx, threshold).Return(users, nil)
		mockGitHub.EXPECT().GetStats(ctx, "user1").Return(stats1, nil)
		mockRepo.EXPECT().SaveUserStats(ctx, repository.Stats{UserID: 1, Commits: 10}).Return(errors.New("save error"))

		err := service.RefreshOutdatedStats(ctx, threshold)
		require.ErrorContains(t, err, "save error")
	})
}
