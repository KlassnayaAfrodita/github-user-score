package kafka

type ScoringRequestMessage struct {
	ApplicationID string `json:"application_id"`
	UserID        int    `json:"user_id"`
	Repos         int    `json:"repos"`
	Stars         int    `json:"stars"`
	Forks         int    `json:"forks"`
	Commits       int    `json:"commits"`
}

type ScoringResultMessage struct {
	ApplicationID string  `json:"application_id"`
	UserID        int     `json:"user_id"`
	Scoring       float64 `json:"scoring"`
}
