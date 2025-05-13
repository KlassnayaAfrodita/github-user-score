package kafka

type ScoringRequestMessage struct {
	ApplicationID int64 `json:"application_id"`
	UserID        int   `json:"user_id"`
	Repos         int   `json:"repos"`
	Stars         int   `json:"stars"`
	Forks         int   `json:"forks"`
	Commits       int   `json:"commits"`
}

type ScoringResultMessage struct {
	ApplicationID int64 `json:"application_id"`
	UserID        int   `json:"user_id"`
	Score         int   `json:"score"`
}
