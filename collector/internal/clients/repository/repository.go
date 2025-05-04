package repository

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/collector/pkg/database"
	"time"
)

type User struct {
	ID       int
	Username string
}

type Stats struct {
	UserID  int
	Repos   int
	Stars   int
	Forks   int
	Commits int
}

type CollectorRepositoryInterface interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, username string) (*User, error)
	SaveUserStats(ctx context.Context, stats Stats) error
}

type CollectorRepository struct{}

func NewCollectorRepository() *CollectorRepositoryInterface {
	return &CollectorRepository{}
}

const SaveStatsQuery = `INSERT INTO user_stats (user_id, repos, stars, forks, commits, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE
		SET repos = EXCLUDED.repos,
		    stars = EXCLUDED.stars,
		    forks = EXCLUDED.forks,
		    commits = EXCLUDED.commits,
		    updated_at = EXCLUDED.updated_at`

func (repo *CollectorRepository) SaveUserStats(ctx context.Context, stats Stats) error {
	tx, err := database.InitTransaction(ctx, "SaveUserStats")
	if err != nil {
		return fmt.Errorf("repository.SaveUserStats: %w", err)
	}

	_, err := tx.Exec(SaveStatsQuery,
		stats.UserID,
		stats.Repos,
		stats.Stars,
		stats.Forks,
		stats.Commits,
		time.Now(),
	)

}

func (repo *CollectorRepository) CreateUser(ctx context.Context, username string) (*User, error) {
	return &User{}, nil
}

func (repo *CollectorRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return &User{}, nil
}
