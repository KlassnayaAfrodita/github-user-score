package main

import (
	"context"
	github "github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/github-api"
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
		log.Fatalf("failed to listen: %v", err)
	}

	kafkaClient := kafka.NewKafkaClient([]string{kafkaBroker}, requestTopic, resultTopic, testGroupID)

	service := services.NewScoringManagerService(repo, collectorClient, kafkaClient)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.Registerscoring - managerServiceServer(grpcServer, controllers.Newscoring-managerHandler(service))

	reflection.Register(grpcServer)

	go func() {
		log.Println("starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	waitForShutdown()
	log.Println("shutting down...")
	grpcServer.GracefulStop()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
