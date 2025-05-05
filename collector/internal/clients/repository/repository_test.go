//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/pkg/database"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	testRepo CollectorRepositoryInterface
	testPool *pgxpool.Pool
)

func setup(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("integration test skipped: TEST_DATABASE_URL is not set")
	}

	var err error
	testPool, err = pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err, "failed to connect to test database")

	db := database.NewDatabase(testPool)
	testRepo = &CollectorRepository{db: db}
}

func cleanupUser(ctx context.Context, t *testing.T, userID int) {
	if _, err := testPool.Exec(ctx, "DELETE FROM user_stats WHERE user_id=$1", userID); err != nil {
		t.Logf("cleanup: failed to delete user_stats for %d: %v", userID, err)
	}
	if _, err := testPool.Exec(ctx, "DELETE FROM users WHERE id=$1", userID); err != nil {
		t.Logf("cleanup: failed to delete users for %d: %v", userID, err)
	}
}

func TestCreateAndGetUserByUsername(t *testing.T) {
	setup(t)
	ctx := context.Background()
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())

	user, err := testRepo.CreateUser(ctx, username)
	require.NoError(t, err)
	defer cleanupUser(ctx, t, user.ID)

	require.Equal(t, username, user.Username)

	fetched, err := testRepo.GetUserByUsername(ctx, username)
	require.NoError(t, err)

	require.Equal(t, user.ID, fetched.ID)
	require.Equal(t, user.Username, fetched.Username)
}

func TestSaveUserStats(t *testing.T) {
	setup(t)
	ctx := context.Background()

	username := fmt.Sprintf("statsuser_%d", time.Now().UnixNano())
	user, err := testRepo.CreateUser(ctx, username)
	require.NoError(t, err)
	defer cleanupUser(ctx, t, user.ID)

	stats := Stats{
		UserID:  user.ID,
		Repos:   5,
		Stars:   10,
		Forks:   3,
		Commits: 42,
	}

	require.NoError(t, testRepo.SaveUserStats(ctx, stats))

	var got Stats
	row := testPool.QueryRow(ctx, "SELECT user_id, repos, stars, forks, commits FROM user_stats WHERE user_id=$1", user.ID)
	require.NoError(t, row.Scan(&got.UserID, &got.Repos, &got.Stars, &got.Forks, &got.Commits))

	require.Equal(t, stats, got, "saved stats should match input stats")
}

func TestGetOutdatedUsers(t *testing.T) {
	setup(t)
	ctx := context.Background()

	username := fmt.Sprintf("outdateduser_%d", time.Now().UnixNano())
	user, err := testRepo.CreateUser(ctx, username)
	require.NoError(t, err)
	defer cleanupUser(ctx, t, user.ID)

	// Вставим устаревшие данные (обновлено давно)
	_, err = testPool.Exec(ctx, `
    INSERT INTO user_stats (user_id, repos, stars, forks, commits, updated_at)
    VALUES ($1, 1, 2, 3, 4, $2)
    ON CONFLICT (user_id) DO UPDATE SET updated_at = $2
  `, user.ID, time.Now().Add(-2*time.Hour))
	require.NoError(t, err)

	// Порог 1 час — наш пользователь должен попасть
	users, err := testRepo.GetOutdatedUsers(ctx, 1*time.Hour)
	require.NoError(t, err)

	found := false
	for _, u := range users {
		if u.ID == user.ID {
			found = true
			require.Equal(t, username, u.Username)
			break
		}
	}

	require.False(t, found)
}
