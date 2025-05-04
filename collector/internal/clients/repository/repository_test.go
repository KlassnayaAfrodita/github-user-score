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
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

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
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	defer cleanupUser(ctx, t, user.ID)

	if user.Username != username {
		t.Errorf("expected username %q, got %q", username, user.Username)
	}

	fetched, err := testRepo.GetUserByUsername(ctx, username)
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if fetched.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, fetched.ID)
	}
	if fetched.Username != user.Username {
		t.Errorf("expected username %q, got %q", user.Username, fetched.Username)
	}
}

func TestSaveUserStats(t *testing.T) {
	setup(t)
	ctx := context.Background()

	username := fmt.Sprintf("statsuser_%d", time.Now().UnixNano())
	user, err := testRepo.CreateUser(ctx, username)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	defer cleanupUser(ctx, t, user.ID)

	stats := Stats{
		UserID:  user.ID,
		Repos:   5,
		Stars:   10,
		Forks:   3,
		Commits: 42,
	}

	err = testRepo.SaveUserStats(ctx, stats)
	if err != nil {
		t.Fatalf("SaveUserStats failed: %v", err)
	}

	var got Stats
	row := testPool.QueryRow(ctx, "SELECT user_id, repos, stars, forks, commits FROM user_stats WHERE user_id=$1", user.ID)
	if err := row.Scan(&got.UserID, &got.Repos, &got.Stars, &got.Forks, &got.Commits); err != nil {
		t.Fatalf("failed to fetch saved stats: %v", err)
	}

	if got != stats {
		t.Errorf("saved stats mismatch, expected %+v, got %+v", stats, got)
	}
}
