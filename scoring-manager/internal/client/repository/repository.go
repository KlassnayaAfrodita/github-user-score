package repository

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/database"
	"github.com/jackc/pgx/v5"
)

type ScoringRepositoryInterface interface {
	CreateScoringApplication(ctx context.Context, app ScoringApplication) (ScoringApplication, error)
	UpdateScoringApplicationStatus(ctx context.Context, appID int64, status ScoringStatus) error
	SaveScoringApplicationResult(ctx context.Context, app ScoringApplication) error
	GetScoringApplicationByID(ctx context.Context, appID int64) (ScoringApplication, error)
	GetExpiredApplications(ctx context.Context, maxAgeMinutes int) ([]int64, error)
	MarkExpiredApplications(ctx context.Context, appIDs []int64) error
}

type ScoringRepository struct {
	db *database.Database
}

func NewScoringRepository(db *database.Database) *ScoringRepository {
	return &ScoringRepository{db: db}
}

const createScoringApplicationQuery = `
INSERT INTO scoring_status (user_id, status)
VALUES ($1, $2)
RETURNING application_id
`

func (repo *ScoringRepository) CreateScoringApplication(ctx context.Context, app ScoringApplication) (ScoringApplication, error) {
	tx, err := repo.db.InitTransaction(ctx, "CreateScoringApplication")
	if err != nil {
		return ScoringApplication{}, fmt.Errorf("repository.CreateScoringApplication: %w", err)
	}
	defer tx.Rollback(ctx)

	var appID int64
	err = tx.QueryRow(ctx, createScoringApplicationQuery, app.UserID, app.Status).Scan(&appID)
	if err != nil {
		return ScoringApplication{}, fmt.Errorf("repository.CreateScoringApplication: %w", err)
	}

	app.ApplicationID = appID

	if err := tx.Commit(ctx); err != nil {
		return ScoringApplication{}, fmt.Errorf("repository.CreateScoringApplication: commit failed: %w", err)
	}

	return app, nil
}

const updateScoringApplicationStatusQuery = `UPDATE scoring_status SET status = $1 WHERE application_id = $2`

func (repo *ScoringRepository) UpdateScoringApplicationStatus(ctx context.Context, appID int64, status ScoringStatus) error {
	tx, err := repo.db.InitTransaction(ctx, "UpdateScoringApplicationStatus")
	if err != nil {
		return fmt.Errorf("repository.UpdateScoringApplicationStatus: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updateScoringApplicationStatusQuery, status, appID)
	if err != nil {
		return fmt.Errorf("repository.UpdateScoringApplicationStatus: %w", err)
	}

	return tx.Commit(ctx)
}

const saveScoringApplicationResultQuery = `INSERT INTO scoring_result (application_id, user_id, score)
		VALUES ($1, $2, $3)`

func (repo *ScoringRepository) SaveScoringApplicationResult(ctx context.Context, app ScoringApplication) error {
	tx, err := repo.db.InitTransaction(ctx, "SaveScoringApplicationResult")
	if err != nil {
		return fmt.Errorf("repository.SaveScoringApplicationResult: %w", err)
	}
	defer tx.Rollback(ctx)

	if app.Score == nil {
		return fmt.Errorf("repository.SaveScoringApplicationResult: score is nil")
	}

	_, err = tx.Exec(ctx, saveScoringApplicationResultQuery, app.ApplicationID, app.UserID, *app.Score)
	if err != nil {
		return fmt.Errorf("repository.SaveScoringApplicationResult: %w", err)
	}

	return tx.Commit(ctx)
}

const getScoringApplicationByID = `SELECT s.application_id, s.user_id, s.status, r.score
		FROM scoring_status s
		LEFT JOIN scoring_result r ON s.application_id = r.application_id
		WHERE s.application_id = $1`

func (repo *ScoringRepository) GetScoringApplicationByID(ctx context.Context, appID int64) (ScoringApplication, error) {
	tx, err := repo.db.InitTransaction(ctx, "UpdateScoringApplicationStatus")
	if err != nil {
		return ScoringApplication{}, fmt.Errorf("repository.GetScoringApplicationByID: %w", err)
	}
	defer tx.Rollback(ctx)

	var app ScoringApplication
	var score *int

	row := tx.QueryRow(ctx, getScoringApplicationByID, appID)
	err = row.Scan(&app.ApplicationID, &app.UserID, &app.Status, &score)
	switch err {
	case nil:
	// продолжаем просто
	case pgx.ErrNoRows:
		return ScoringApplication{}, nil
	default:
		return ScoringApplication{}, fmt.Errorf("repository.GetUserByUsername: %w", err)
	}

	app.Score = score

	return app, tx.Commit(ctx)
}

const getExpiredApplicationsQuery = `
	SELECT application_id
	FROM scoring_status
	WHERE status = $1 AND created_at < NOW() - ($2 || ' minutes')::interval;
`

func (repo *ScoringRepository) GetExpiredApplications(ctx context.Context, maxAgeMinutes int) ([]int64, error) {
	tx, err := repo.db.InitTransaction(ctx, "GetExpiredApplications")
	if err != nil {
		return nil, fmt.Errorf("repository.GetExpiredApplications: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, getExpiredApplicationsQuery, StatusInitial, fmt.Sprint(maxAgeMinutes))
	if err != nil {
		return nil, fmt.Errorf("repository.GetExpiredApplications: %w", err)
	}
	defer rows.Close()

	var appIDs []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		switch err {
		case nil:
		// продолжаем просто
		case pgx.ErrNoRows:
			return []int64{}, nil
		default:
			return []int64{}, fmt.Errorf("repository.GetExpiredApplications: %w", err)
		}
		appIDs = append(appIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.GetExpiredApplications: %w", err)
	}

	return appIDs, tx.Commit(ctx)
}

const markApplicationsAsFailedQuery = `
	UPDATE scoring_status
	SET status = $1
	WHERE application_id = ANY($2)
`

func (repo *ScoringRepository) MarkExpiredApplications(ctx context.Context, appIDs []int64) error {
	if len(appIDs) == 0 {
		return nil
	}

	tx, err := repo.db.InitTransaction(ctx, "MarkExpiredApplications")
	if err != nil {
		return fmt.Errorf("repository.MarkExpiredApplications: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, markApplicationsAsFailedQuery, StatusFailed, appIDs)
	if err != nil {
		return fmt.Errorf("repository.MarkExpiredApplications: %w", err)
	}
	return tx.Commit(ctx)
}
