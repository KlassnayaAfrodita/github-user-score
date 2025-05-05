//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/KlassnayaAfrodita/github-user-score/collector/pkg/database"
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
