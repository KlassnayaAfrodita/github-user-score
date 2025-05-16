//go:build integration
// +build integration

package grpc

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/pkg/pb/scoring-manager/api"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestScoringManagerClient_FullIntegration(t *testing.T) {
	address := os.Getenv("SCORING_MANAGER_ADDRESS")
	require.NotEmpty(t, address, "SCORING_MANAGER_ADDRESS must be set")

	client, err := NewScoringManagerClient(address, 5*time.Second)
	require.NoError(t, err)

	// 1. Стартуем скоринг
	startResp, err := client.(*ScoringManagerClient).service.StartScoring(context.Background(), &api.StartScoringRequest{
		Username: "some-existing-github-username",
	})
	require.NoError(t, err)
	applicationID := startResp.GetApplicationId()

	// 2. Poll GetStatus до SUCCESS / FAILED или timeout
	var result GetScoreResponse
	const timeout = 60 * time.Second
	const interval = 2 * time.Second
	deadline := time.Now().Add(timeout)

	for {
		result, err = client.GetStatus(context.Background(), applicationID)
		if err == nil && (result.Status == StatusSuccess || result.Status == StatusFailed) {
			break
		}

		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for scoring to complete")
		}

		time.Sleep(interval)
	}

	t.Logf("Final scoring result: status=%v, score=%d", result.Status, result.Scoring)
	require.Equal(t, StatusSuccess, result.Status)
}
