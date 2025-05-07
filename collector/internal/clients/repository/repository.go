package repository

import (
	"context"
	"fmt"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/pkg/database"
	"github.com/jackc/pgx/v5"
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
	GetOutdatedUsers(ctx context.Context, threshold time.Duration) ([]User, error)
	GetUserStats(ctx context.Context, userID int) (Stats, error)
}

type CollectorRepository struct {
	db *database.Database
}

func NewCollectorRepository(db *database.Database) *CollectorRepository {
	return &CollectorRepository{db: db}
}

const saveStatsQuery = `INSERT INTO user_stats (user_id, repos, stars, forks, commits, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE
		SET repos = EXCLUDED.repos,
		    stars = EXCLUDED.stars,
		    forks = EXCLUDED.forks,
		    commits = EXCLUDED.commits,
		    updated_at = EXCLUDED.updated_at`

func (repo *CollectorRepository) SaveUserStats(ctx context.Context, stats Stats) error {
	tx, err := repo.db.InitTransaction(ctx, "SaveUserStats")
	if err != nil {
		return fmt.Errorf("repository.SaveUserStats: %w", err)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, saveStatsQuery,
		stats.UserID,
		stats.Repos,
		stats.Stars,
		stats.Forks,
		stats.Commits,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("repository.SaveUserStats: %w", err)
	}

	tx.Commit(ctx)
	return nil
}

const createUser = `INSERT INTO users (username) VALUES ($1) RETURNING id`

func (repo *CollectorRepository) CreateUser(ctx context.Context, username string) (*User, error) {
	tx, err := repo.db.InitTransaction(ctx, "CreateUser")
	if err != nil {
		return &User{}, fmt.Errorf("repository.CreateUser: %w", err)
	}

	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx, createUser, username).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("repository.CreateUser: %w", err)
	}

	tx.Commit(ctx)
	return &User{
		ID:       id,
		Username: username,
	}, nil
}

const getByUsername = `SELECT id, username FROM users WHERE username = $1`

func (repo *CollectorRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	tx, err := repo.db.InitTransaction(ctx, "GetUserByUsername")
	if err != nil {
		return nil, fmt.Errorf("repository.GetUserByUsername: %w", err)
	}

	defer tx.Rollback(ctx)

	var user User
	row := tx.QueryRow(ctx, getByUsername, username)

	err = row.Scan(&user.ID, &user.Username)
	switch err {
	case nil:
	// продолжаем просто
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, fmt.Errorf("repository.GetUserByUsername: %w", err)
	}

	tx.Commit(ctx)
	return &user, nil
}

const getOutdatedUsersQuery = `
  SELECT id, username FROM users 
  JOIN user_stats ON users.id = user_stats.user_id 
  WHERE user_stats.updated_at < NOW() - make_interval(secs => $1)
`

func (repo *CollectorRepository) GetOutdatedUsers(ctx context.Context, threshold time.Duration) ([]User, error) {
	tx, err := repo.db.InitTransaction(ctx, "GetOutdatedUsers")
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	seconds := int64(threshold.Seconds())

	rows, err := tx.Query(ctx, getOutdatedUsersQuery, seconds)
	if err != nil {
		return nil, fmt.Errorf("repository.GetUserByUsername: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username)
		switch err {
		case nil:
		// продолжаем просто
		case pgx.ErrNoRows:
			return nil, nil
		default:
			return nil, fmt.Errorf("repository.GetUserByUsername: %w", err)
		}
		users = append(users, user)
	}

	tx.Commit(ctx)
	return users, nil
}

const getUserStats = `SELECT user_id, repos, stars, forks, commits FROM user_stats WHERE user_id = $1`

func (repo *CollectorRepository) GetUserStats(ctx context.Context, userID int) (Stats, error) {
	tx, err := repo.db.InitTransaction(ctx, "GetUserStats")
	if err != nil {
		return Stats{}, fmt.Errorf("repository.GetUserStats: %w", err)
	}
	defer tx.Rollback(ctx)

	var stats Stats
	err = tx.QueryRow(ctx, getUserStats, userID).Scan(
		&stats.UserID,
		&stats.Repos,
		&stats.Stars,
		&stats.Forks,
		&stats.Commits,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("repository.GetUserStats: %w", err)
	}

	tx.Commit(ctx)
	return stats, nil
}
