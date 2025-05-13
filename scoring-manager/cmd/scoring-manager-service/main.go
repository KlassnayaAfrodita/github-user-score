package main

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/pkg/database"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/services"
	pb "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/pkg/pb/scoring-manager/api"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	collector "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/grpc"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/kafka"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/controllers"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	dbURL := "postgres://testuser:testpass@localhost:5434/test_db_scoring-manager?sslmode=disable"

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	db := database.NewDatabase(pool)
	repo := repository.NewScoringRepository(db)

	collectorClient, err := collector.NewCollectorClient("localhost:50051", 5*time.Second)
	if err != nil {
		log.Fatalf("failed to connect to collector: %v", err)
	}

	kafkaClient := kafka.NewKafkaClient([]string{kafkaBroker}, requestTopic, resultTopic, testGroupID)

	service := services.NewScoringManagerService(repo, collectorClient, kafkaClient)

	controller := controllers.NewScoringManagerHandler(service)

	server := grpc.NewServer()
	pb.RegisterScoringManagerServiceServer(server, controller)
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("starting Kafka consumer loop...")
		if err = service.ConsumingScoringResults(ctx); err != nil {
			log.Printf("Kafka consumer stopped with error: %v", err)
		}
	}()

	go func() {
		log.Println("starting gRPC server on :50052")
		if err = server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	waitForShutdown()
	log.Println("shutting down...")
	server.GracefulStop()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
