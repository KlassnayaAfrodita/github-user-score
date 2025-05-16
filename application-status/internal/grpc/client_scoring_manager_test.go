//go:build integration
// +build integration

package grpc

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestScoringManagerClient_FullIntegration(t *testing.T) {
	address := os.Getenv("SCORING_MANAGER_ADDRESS")

	clientRaw, err := NewScoringManagerClient(address, 5*time.Second)
	assert.NoError(t, err)
	client := clientRaw.(*ScoringManagerClient)

	// Запускаем скоринг
	username := "testuser123"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startResp, err := client.service.StartScoring(ctx, &grpc.StartScoringRequest{
		Username: username,
	})
	assert.NoError(t, err)
	assert.Greater(t, startResp.ApplicationId, int64(0))

	applicationID := startResp.ApplicationId

	// Опрашиваем GetStatus с интервалом
	var (
		finalStatus grpc.ScoringStatus
		score       int32
	)

	maxWait := time.Now().Add(100 * time.Second)
	for {
		resp, err := client.GetStatus(context.Background(), applicationID)
		assert.NoError(t, err)

		if resp.Status == grpc.StatusSuccess || resp.Status == grpc.StatusFailed {
			finalStatus = resp.Status
			score = resp.Scoring
			break
		}

		if time.Now().After(maxWait) {
			t.Fatal("Timeout waiting for scoring to complete")
		}
		time.Sleep(2 * time.Second)
	}

	t.Logf("Final scoring status: %v, score: %d", finalStatus, score)

	// Проверка результата
	assert.Equal(t, grpc.StatusSuccess, finalStatus, "Expected successful scoring")
	assert.GreaterOrEqual(t, score, int32(0))
	assert.LessOrEqual(t, score, int32(100))
}
