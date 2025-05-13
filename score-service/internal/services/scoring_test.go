package service

import (
	"context"
	"github.com/KlassnayaAfrodita/github-user-score/score-service/internal/client/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/KlassnayaAfrodita/github-user-score/score-service/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
)

func TestScoringService_handleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafkaClient := mocks.NewMockKafkaClientInterface(ctrl)
	service := NewScoringService(mockKafkaClient)

	msg := kafka.ScoringRequestMessage{
		ApplicationID: 1,
		UserID:        42,
		Repos:         10,
		Stars:         50,
		Forks:         5,
		Commits:       100,
	}

	mockKafkaClient.EXPECT().
		ProduceScoringResult(gomock.Any(), gomock.AssignableToTypeOf(kafka.ScoringResultMessage{})).
		DoAndReturn(func(ctx context.Context, result kafka.ScoringResultMessage) error {
			assert.Equal(t, msg.ApplicationID, result.ApplicationID)
			assert.Equal(t, msg.UserID, result.UserID)
			assert.True(t, result.Score >= 0 && result.Score <= 100)
			return nil
		})

	err := service.HandleMessage(msg)
	require.NoError(t, err)
}
