package service

import "github.com/KlassnayaAfrodita/github-user-score/score-service/internal/client/kafka"

type ScoringService struct {
	kafka kafka.KafkaClientInterface
}

func NewScoringService(kc kafka.KafkaClientInterface) *ScoringService {
	return &ScoringService{kafka: kc}
}
