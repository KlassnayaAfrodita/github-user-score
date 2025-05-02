package repository

import "context"

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

//const SaveStatsQuery = `insert into `

func (repo *CollectorRepository) SaveUserStats(ctx context.Context, stats Stats) error {
	return nil
}

func (repo *CollectorRepository) CreateUser(ctx context.Context, username string) (*User, error) {
	return &User{}, nil
}

func (repo *CollectorRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return &User{}, nil
}
