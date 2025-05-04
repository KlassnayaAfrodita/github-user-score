package github_api

type Stats struct {
	Repos   int
	Stars   int
	Forks   int
	Commits int
}

type GitHubClient interface {
	GetStats(username string) (Stats, error)
}

// пока заглушка
func GetStats(username string) (Stats, error) {
	return Stats{
		Repos:   1,
		Stars:   1,
		Forks:   1,
		Commits: 1,
	}, nil
}
