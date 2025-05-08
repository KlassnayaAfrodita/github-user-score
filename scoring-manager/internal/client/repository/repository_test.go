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

	app := &ScoringApplication{
		UserID: 123456,
		Status: StatusInitial,
	}

	err := scoringRepo.CreateScoringApplication(ctx, app)
	require.NoError(t, err)
	require.NotZero(t, app.ApplicationID)

	defer cleanupScoring(ctx, t, app.ApplicationID)

	fetched, err := scoringRepo.GetScoringApplicationByID(ctx, fmt.Sprint(app.ApplicationID))
	require.NoError(t, err)

	require.Equal(t, app.UserID, fetched.UserID)
	require.Equal(t, app.Status, fetched.Status)
	require.Nil(t, fetched.Score)
}

func TestUpdateScoringApplicationStatus(t *testing.T) {
	setup(t)
	ctx := context.Background()

	app := &ScoringApplication{
		UserID: 222222,
		Status: StatusInitial,
	}
	require.NoError(t, scoringRepo.CreateScoringApplication(ctx, app))
	defer cleanupScoring(ctx, t, app.ApplicationID)

	err := scoringRepo.UpdateScoringApplicationStatus(ctx, app.ApplicationID, StatusSuccess)
	require.NoError(t, err)

	fetched, err := scoringRepo.GetScoringApplicationByID(ctx, fmt.Sprint(app.ApplicationID))
	require.NoError(t, err)
	require.Equal(t, StatusSuccess, fetched.Status)
}

func TestSaveScoringApplicationResult(t *testing.T) {
	setup(t)
	ctx := context.Background()

	app := &ScoringApplication{
		UserID: 333333,
		Status: StatusInitial,
	}
	require.NoError(t, scoringRepo.CreateScoringApplication(ctx, app))
	defer cleanupScoring(ctx, t, app.ApplicationID)

	score := 95
	app.Score = &score
	require.NoError(t, scoringRepo.SaveScoringApplicationResult(ctx, *app))

	fetched, err := scoringRepo.GetScoringApplicationByID(ctx, fmt.Sprint(app.ApplicationID))
	require.NoError(t, err)
	require.NotNil(t, fetched.Score)
	require.Equal(t, score, *fetched.Score)
}
