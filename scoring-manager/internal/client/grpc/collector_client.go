package grpc_client

import (
	"context"
	"fmt"
	collector "github.com/KlassnayaAfrodita/github-user-score/collector/pkg/pb/collector/api"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/utils"
	"google.golang.org/grpc"
	"time"
)

type CollectorClientInterface interface {
	GetUserStats(ctx context.Context, username string) (repository.UserStats, error)
}

type CollectorClient struct {
	conn    *grpc.ClientConn
	service collector.CollectorServiceClient
}

func NewCollectorClient(address string, timeout time.Duration) (CollectorClientInterface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &CollectorClient{
		conn:    conn,
		service: collector.NewCollectorServiceClient(conn),
	}, nil
}

func (cc *CollectorClient) GetUserStats(ctx context.Context, username string) (repository.UserStats, error) {
	req := &collector.CollectUserStatsRequest{Username: username}
	userStats, err := cc.service.CollectUserStats(ctx, req)
	if err != nil {
		return repository.UserStats{}, fmt.Errorf("CollectorClient.GetUserStats: %w", err)
	}

	return utils.ToUserStats(userStats), nilss
}
