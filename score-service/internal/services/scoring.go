package service

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/KlassnayaAfrodita/github-user-score/score-service/internal/client/kafka"
)

func (s *ScoringService) Start(ctx context.Context) error {
	return s.kafka.ConsumeScoringRequests(ctx, s.HandleMessage)
}

func (s *ScoringService) HandleMessage(msg kafka.ScoringRequestMessage) error {
	ctx := context.Background()
	time.Sleep(time.Duration(rand.Intn(61)) * time.Second)
	rand.Seed(time.Now().UnixNano())

	result := kafka.ScoringResultMessage{
		ApplicationID: msg.ApplicationID,
		UserID:        msg.UserID,
		Score:         rand.Intn(101),
	}

	log.Printf("Processed scoring: %+v -> %+v", msg, result)
	return s.kafka.ProduceScoringResult(ctx, result)
}
