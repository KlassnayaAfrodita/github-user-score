//go:build integration
// +build integration

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/KlassnayaAfrodita/github-user-score/scoring_manager/internal/clients/kafka"
	"github.com/stretchr/testify/require"
)

const (
	kafkaBroker  = "localhost:9092"
	requestTopic = "scoring_requests"
	resultTopic  = "scoring_results"
	groupID      = "test-group"
)

func TestKafkaClient_PublishAndConsume(t *testing.T) {
	// Подготовка клиента
	client := kafka.NewKafkaClient([]string{kafkaBroker}, requestTopic, resultTopic, groupID)

	ctx := context.Background()

	// Пример сообщения
	msg := kafka.ScoringResultMessage{
		ApplicationID: fmt.Sprintf("app-%d", time.Now().UnixNano()),
		UserID:        123,
		Scoring:       87.5,
	}

	// Слушаем resultTopic в горутине
	done := make(chan bool, 1)
	go func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err := client.ConsumeScoringResults(ctx, func(m kafka.ScoringResultMessage) error {
			require.Equal(t, msg.ApplicationID, m.ApplicationID)
			require.Equal(t, msg.UserID, m.UserID)
			require.Equal(t, msg.Scoring, m.Scoring)
			done <- true
			return nil
		})
		require.NoError(t, err)
	}()

	// Даем Kafka немного времени "проснуться"
	time.Sleep(2 * time.Second)

	// Публикуем результат
	raw, err := json.Marshal(msg)
	require.NoError(t, err)

	writer := client // или использовать прямой вызов
	err = writer.Producer().WriteMessages(ctx, kafkaMessage(msg.ApplicationID, raw))
	require.NoError(t, err)

	select {
	case <-done:
		// Всё прошло успешно
	case <-time.After(10 * time.Second):
		t.Fatal("did not receive scoring result in time")
	}
}

func kafkaMessage(key string, value []byte) kafka.Message {
	return kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
}
