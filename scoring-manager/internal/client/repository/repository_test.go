//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	scoringRepo ScoringRepositoryInterface
	pool        *pgxpool.Pool
)

func setup(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	db := database.NewDatabase(pool)
	scoringRepo = NewScoringRepository(db)
}

func cleanupScoring(ctx context.Context, t *testing.T, appID int64) {
	_, _ = pool.Exec(ctx, "DELETE FROM scoring_result WHERE application_id = $1", appID)
	_, _ = pool.Exec(ctx, "DELETE FROM scoring_status WHERE application_id = $1", appID)
}

func TestCreateAndGetScoringApplication(t *testing.T) {
	setup(t)
	ctx := context.Background()

	app := ScoringApplication{
		UserID: 123456,
		Status: StatusInitial,
	}

	createdApp, err := scoringRepo.CreateScoringApplication(ctx, app)
	require.NoError(t, err)
	require.NotZero(t, createdApp.ApplicationID)

	defer cleanupScoring(ctx, t, createdApp.ApplicationID)

	fetched, err := scoringRepo.GetScoringApplicationByID(ctx, fmt.Sprint(createdApp.ApplicationID))
	require.NoError(t, err)

	require.Equal(t, createdApp.UserID, fetched.UserID)
	require.Equal(t, createdApp.Status, fetched.Status)
	require.Nil(t, fetched.Score)
}

func TestUpdateScoringApplicationStatus(t *testing.T) {
	setup(t)
	ctx := context.Background()

	app := ScoringApplication{
		UserID: 222222,
		Status: StatusInitial,
	}
	createdApp, err := scoringRepo.CreateScoringApplication(ctx, app)
	require.NoError(t, err)
	defer cleanupScoring(ctx, t, createdApp.ApplicationID)

	err = scoringRepo.UpdateScoringApplicationStatus(ctx, createdApp.ApplicationID, StatusSuccess)
	require.NoError(t, err)

	fetched, err := scoringRepo.GetScoringApplicationByID(ctx, fmt.Sprint(createdApp.ApplicationID))
	require.NoError(t, err)
	require.Equal(t, StatusSuccess, fetched.Status)
}

func TestSaveScoringApplicationResult(t *testing.T) {
	setup(t)
	ctx := context.Background()

	app := ScoringApplication{
		UserID: 333333,
		Status: StatusInitial,
	}
	createdApp, err := scoringRepo.CreateScoringApplication(ctx, app)
	require.NoError(t, err)
	defer cleanupScoring(ctx, t, createdApp.ApplicationID)

	score := 95
	createdApp.Score = &score
	require.NoError(t, scoringRepo.SaveScoringApplicationResult(ctx, createdApp))

	fetched, err := scoringRepo.GetScoringApplicationByID(ctx, fmt.Sprint(createdApp.ApplicationID))
	require.NoError(t, err)
	require.NotNil(t, fetched.Score)
	require.Equal(t, score, *fetched.Score)
}

func forceInitialOldStatus(ctx context.Context, t *testing.T, appID int64, minutesAgo int) {
	_, err := pool.Exec(ctx, `
		UPDATE scoring_status
		SET created_at = NOW() - ($1 || ' minutes')::interval
		WHERE application_id = $2`, fmt.Sprint(minutesAgo), appID)
	require.NoError(t, err)
}

func TestGetExpiredApplications(t *testing.T) {
	setup(t)
	ctx := context.Background()

	app := ScoringApplication{
		UserID: 444444,
		Status: StatusInitial,
	}
	createdApp, err := scoringRepo.CreateScoringApplication(ctx, app)
	require.NoError(t, err)
	defer cleanupScoring(ctx, t, createdApp.ApplicationID)

	forceInitialOldStatus(ctx, t, createdApp.ApplicationID, 10)

	expiredIDs, err := scoringRepo.GetExpiredApplications(ctx, 5)
	require.NoError(t, err)

	require.Contains(t, expiredIDs, createdApp.ApplicationID)
}

func TestMarkExpiredApplications(t *testing.T) {
	setup(t)
	ctx := context.Background()

	app := ScoringApplication{
		UserID: 555555,
		Status: StatusInitial,
	}
	createdApp, err := scoringRepo.CreateScoringApplication(ctx, app)
	require.NoError(t, err)
	defer cleanupScoring(ctx, t, createdApp.ApplicationID)

	forceInitialOldStatus(ctx, t, createdApp.ApplicationID, 10)

	err = scoringRepo.MarkExpiredApplications(ctx, []int64{createdApp.ApplicationID})
	require.NoError(t, err)

	fetched, err := scoringRepo.GetScoringApplicationByID(ctx, fmt.Sprint(createdApp.ApplicationID))
	require.NoError(t, err)
	require.Equal(t, StatusFailed, fetched.Status)
}
