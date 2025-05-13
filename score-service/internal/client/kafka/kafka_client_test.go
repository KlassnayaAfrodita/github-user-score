//go:build integration
// +build integration

package kafka

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

// Местная реализация handle-функции (вместо импорта services)
func handleMessage(ctx context.Context, kc *KafkaClient) func(ScoringRequestMessage) error {
	return func(msg ScoringRequestMessage) error {
		time.Sleep(time.Duration(rand.Intn(61)) * time.Second)
		rand.Seed(time.Now().UnixNano())

		result := ScoringResultMessage{
			ApplicationID: msg.ApplicationID,
			UserID:        msg.UserID,
			Score:         rand.Intn(101),
		}

		log.Printf("Integration processed: %+v -> %+v", msg, result)
		return kc.ProduceScoringResult(ctx, result)
	}
}

func TestKafkaClient_Integration(t *testing.T) {
	brokers := []string{"localhost:9092"}
	requestTopic := "scoring_requests"
	resultTopic := "scoring_results"
	groupID := "scoring_integration"

	kc := NewKafkaClient(brokers, requestTopic, resultTopic, groupID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем consume в фоне
	go func() {
		err := kc.ConsumeScoringRequests(ctx, handleMessage(ctx, kc))
		if err != nil {
			log.Printf("consume error: %v", err)
		}
	}()

	// Отправляем сообщение
	producer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    requestTopic,
		Balancer: &kafka.LeastBytes{},
	}

	req := ScoringRequestMessage{
		ApplicationID: 123,
		UserID:        456,
		Repos:         10,
		Stars:         20,
		Forks:         5,
		Commits:       100,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	err = producer.WriteMessages(ctx, kafka.Message{Value: data})
	require.NoError(t, err)

	// Чтение результата
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   resultTopic,
		GroupID: groupID + "-results",
	})
	defer consumer.Close()

	done := make(chan ScoringResultMessage)

	go func() {
		for {
			m, err := consumer.ReadMessage(ctx)
			if err != nil {
				log.Printf("read error: %v", err)
				continue
			}

			var res ScoringResultMessage
			if err := json.Unmarshal(m.Value, &res); err == nil && res.ApplicationID == req.ApplicationID {
				done <- res
				return
			}
		}
	}()

	select {
	case res := <-done:
		require.Equal(t, req.UserID, res.UserID)
		require.True(t, res.Score >= 0 && res.Score <= 100)
	case <-time.After(100 * time.Second):
		t.Fatal("timeout waiting for scoring result")
	}
}
