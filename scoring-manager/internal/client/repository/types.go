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
	ApplicationID int64
	UserID        int
	Status        ScoringStatus
	Score         *int // указатель, так как score может быть null
}
