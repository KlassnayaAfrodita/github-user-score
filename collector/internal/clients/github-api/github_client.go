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
