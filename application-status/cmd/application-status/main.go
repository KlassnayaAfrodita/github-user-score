package application_status

import (
	"context"
	client "github.com/KlassnayaAfrodita/github-user-score/application-status/internal/client/grpc"
	"github.com/KlassnayaAfrodita/github-user-score/application-status/internal/controllers"
	"github.com/KlassnayaAfrodita/github-user-score/application-status/internal/services"
	api "github.com/KlassnayaAfrodita/github-user-score/application-status/pkg/pb/application-status/api"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	scoringManageClient, err := client.NewScoringManagerClient("scoring-manager:50052", 5*time.Second)
	if err != nil {
		log.Fatalf("failed to connect to scoring-manager: %v", err)
	}

	service := services.NewApplicationStatusClient(scoringManageClient)

	controller := controllers.NewApplicationStatusControlles(service)

	server := grpc.NewServer()
	api.RegisterApplicationStatusServiceServer(server, controller)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("starting gRPC server on :50053")
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
