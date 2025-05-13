//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	pb "github.com/KlassnayaAfrodita/github-user-score/api/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestEndToEndScoring(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // gRPC scoring-manager
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewScoringManagerClient(conn)

	ctx := context.Background()

	// 1. Отправляем StartScoring
	resp, err := client.StartScoring(ctx, &pb.StartScoringRequest{
		Username: "torvalds", // например
	})
	require.NoError(t, err)
	require.NotZero(t, resp.ApplicationId)

	appID := resp.ApplicationId

	// 2. Ожидаем результата
	var result *pb.GetResultResponse
	for i := 0; i < 10; i++ {
		status, _ := client.GetStatus(ctx, &pb.GetStatusRequest{ApplicationId: appID})
		if status != nil && status.Status == "success" {
			result, err = client.GetResult(ctx, &pb.GetResultRequest{ApplicationId: appID})
			require.NoError(t, err)
			break
		}
		time.Sleep(20 * time.Second)
	}

	require.NotNil(t, result)
	t.Logf("ApplicationID: %d, Score: %d", result.ApplicationId, result.Score)
}
