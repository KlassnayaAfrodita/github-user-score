package repository

type ScoringStatus int

const (
	//INITIAL
	StatusInitial ScoringStatus = 0
	//SUCCESS
	StatusSuccess ScoringStatus = 1
	//FAILED
	StatusFailed ScoringStatus = 2
)

type ScoringApplication struct {
	ApplicationID string
	UserID        int
	Username      string
	Status        ScoringStatus
	Score         int
}
