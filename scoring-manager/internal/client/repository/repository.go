package repository

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/database"
	"github.com/jackc/pgx/v5"
)

type ScoringRepositoryInterface interface {
	CreateScoringApplication(ctx context.Context, app ScoringApplication) error
	UpdateScoringApplicationStatus(ctx context.Context, appID string, status ScoringStatus) error
	SaveScoringApplicationResult(ctx context.Context, app ScoringApplication) error
	GetScoringApplicationByID(ctx context.Context, appID string) (ScoringApplication, error)
}

type ScoringRepository struct {
	db *database.Database
}

func NewScoringRepository(db *database.Database) *ScoringRepository {
	return &ScoringRepository{db: db}
}

const createScoringApplicationQuery = `INSERT INTO scoring_status (application_id, user_id, status)
VALUES ($1, $2, $3)`

func (repo *ScoringRepository) CreateScoringApplication(ctx context.Context, app ScoringApplication) error {
	tx, err := repo.db.InitTransaction(ctx, "CreateScoringApplication")
	if err != nil {
		return fmt.Errorf("repository.CreateScoringApplication: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, createScoringApplicationQuery, app.ApplicationID, app.UserID, app.Status)
	if err != nil {
		return fmt.Errorf("repository.CreateScoringApplication: %w", err)
	}

	tx.Commit(ctx)
	return nil
}

const updateScoringApplicationStatusQuery = `UPDATE scoring_status SET status = $1 WHERE application_id = $2`

func (repo *ScoringRepository) UpdateScoringApplicationStatus(ctx context.Context, appID string, status ScoringStatus) error {
	tx, err := repo.db.InitTransaction(ctx, "UpdateScoringApplicationStatus")
	if err != nil {
		return fmt.Errorf("repository.UpdateScoringApplicationStatus: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updateScoringApplicationStatusQuery, status, appID)
	if err != nil {
		return fmt.Errorf("repository.UpdateScoringApplicationStatus: %w", err)
	}

	tx.Commit(ctx)
	return nil
}

const saveScoringApplicationResultQuery = `INSERT INTO scoring_result (application_id, user_id, score)
		VALUES ($1, $2, $3)`

func (repo *ScoringRepository) SaveScoringApplicationResult(ctx context.Context, app ScoringApplication) error {
	tx, err := repo.db.InitTransaction(ctx, "UpdateScoringApplicationStatus")
	if err != nil {
		return fmt.Errorf("repository.SaveScoringApplicationResult: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, saveScoringApplicationResultQuery, app.ApplicationID, app.UserID, app.Score)
	if err != nil {
		return fmt.Errorf("repository.SaveScoringApplicationResult: %w", err)
	}

	tx.Commit(ctx)
	return nil
}

const getScoringApplicationByID = `SELECT s.application_id, s.user_id, s.status, r.score
		FROM scoring_status s
		LEFT JOIN scoring_result r ON s.application_id = r.application_id
		WHERE s.application_id = $1`

func (repo *ScoringRepository) GetScoringApplicationByID(ctx context.Context, appID string) (ScoringApplication, error) {
	tx, err := repo.db.InitTransaction(ctx, "UpdateScoringApplicationStatus")
	if err != nil {
		return ScoringApplication{}, fmt.Errorf("repository.GetScoringApplicationByID: %w", err)
	}
	defer tx.Rollback(ctx)

	var app ScoringApplication

	row := tx.QueryRow(ctx, getScoringApplicationByID, appID)
	err = row.Scan(&app.ApplicationID, &app.UserID, &app.Status, &app.Score)
	switch err {
	case nil:
	// продолжаем просто
	case pgx.ErrNoRows:
		return ScoringApplication{}, nil
	default:
		return ScoringApplication{}, fmt.Errorf("repository.GetUserByUsername: %w", err)
	}

	tx.Commit(ctx)
	return app, nil
}
