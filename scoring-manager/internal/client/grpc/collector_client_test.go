//go:build integration
// +build integration

package grpc_client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCollectorClient_GetUserStats(t *testing.T) {
	address := os.Getenv("COLLECTOR_GRPC_ADDR")
	if address == "" {
		t.Skip("COLLECTOR_GRPC_ADDR not set")
	}

	client, err := NewCollectorClient(address, 5*time.Second)
	require.NoError(t, err)
	defer func() {
		if c, ok := client.(*CollectorClient); ok {
			_ = c.conn.Close()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	username := "123"
	stats, err := client.GetUserStats(ctx, username)
	require.NoError(t, err)

	require.NotZero(t, stats.UserID)
	require.Equal(t, stats.Repos, int32(1))   // так как заглушка
	require.Equal(t, stats.Stars, int32(1))   // так как заглушка
	require.Equal(t, stats.Forks, int32(1))   // так как заглушка
	require.Equal(t, stats.Commits, int32(1)) // так как заглушка
}
