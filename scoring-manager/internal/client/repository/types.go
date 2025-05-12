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
	UserID        int32
	Status        ScoringStatus
	Score         *int // указатель, так как score может быть null
}

type UserStats struct {
	UserID  int32
	Repos   int32
	Stars   int32
	Forks   int32
	Commits int32
}
