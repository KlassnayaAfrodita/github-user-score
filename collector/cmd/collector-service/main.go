package main

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/pkg/database"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/services"
	pb "github.com/KlassnayaAfrodita/github-user-score/collector/pkg/pb/collector/api"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	github "github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/github-api"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/controllers"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//dbURL := os.Getenv("DATABASE_URL")
	//if dbURL == "" {
	//	log.Fatal("DATABASE_URL is not set")
	//}

	dbURL := "postgres://testuser:testpass@localhost:5432/test_db_collector?sslmode=disable"

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	db := database.NewDatabase(pool)
	repo := repository.NewCollectorRepository(db)

	githubClient := github.NewGitHubClient()
	service := services.NewCollectorService(repo, githubClient)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCollectorServiceServer(grpcServer, controllers.NewCollectorHandler(service))

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
