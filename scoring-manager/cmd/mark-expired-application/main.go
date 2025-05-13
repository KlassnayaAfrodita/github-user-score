package main

import (
	"context"
	collector "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/grpc"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/database"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	kafkaBroker  = "localhost:9092"
	requestTopic = "scoring-requests"
	resultTopic  = "scoring-results"
	testGroupID  = "test-group"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://testuser:testpass@localhost:5434/test_db_collector?sslmodedisable")
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	db := database.NewDatabase(pool)
	repo := repository.NewScoringRepository(db)
	collectorClient, err := collector.NewCollectorClient("localhost:50051", 5*time.Second)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	kafkaClient := kafka.NewKafkaClient([]string{kafkaBroker}, requestTopic, resultTopic, testGroupID)

	scoringManagerService := services.NewScoringManagerService(repo, collectorClient, kafkaClient)

	startCronJob(ctx, scoringManagerService)

	select {}
}

func startCronJob(ctx context.Context, svc services.ScoringManagerServiceInterface) {
	c := cron.New()

	_, err := c.AddFunc("@every 30m", func() {
		log.Println("[cron] starting outdated stats refresh")
		if err := svc.MarkExpiredApplications(ctx, 15); err != nil {
			log.Printf("[cron] refresh error: %v\n", err)
		} else {
			log.Println("[cron] refresh completed")
		}
	})

	if err != nil {
		log.Fatalf("failed to schedule cron job: %v", err)
	}

	c.Start()
}
