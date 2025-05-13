package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	kafkaclient "github.com/KlassnayaAfrodita/github-user-score/score-service/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/score-service/internal/services"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	kafkaClient := kafkaclient.NewKafkaClient(
		[]string{"localhost:9092"},
		"scoring-requests",
		"scoring-results",
		"scorer-group",
	)

	scoringService := service.NewScoringService(kafkaClient)

	log.Println("Starting scorer service...")
	if err := scoringService.Start(ctx); err != nil {
		log.Fatalf("service error: %v", err)
	}
}
