//go:build integration
// +build integration

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

const (
	kafkaBroker  = "localhost:9092"
	requestTopic = "scoring-requests"
	resultTopic  = "scoring-results"
	testGroupID  = "test-group"
)

func setupKafkaTopic(topic string) error {
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}

func TestKafkaClient_PublishAndConsume(t *testing.T) {
	require.NoError(t, setupKafkaTopic(requestTopic))
	require.NoError(t, setupKafkaTopic(resultTopic))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := NewKafkaClient([]string{kafkaBroker}, requestTopic, resultTopic, testGroupID)

	expected := ScoringResultMessage{
		ApplicationID: int64(42),
		UserID:        123,
		Score:         99,
	}
	data, err := json.Marshal(expected)
	require.NoError(t, err)

	go func() {
		writer := kafka.Writer{
			Addr:     kafka.TCP(kafkaBroker),
			Topic:    resultTopic,
			Balancer: &kafka.LeastBytes{},
		}
		defer writer.Close()

		err := writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(strconv.Itoa(int(expected.ApplicationID))),
			Value: data,
		})
		require.NoError(t, err)
	}()

	err = client.ConsumeScoringResults(ctx, func(received ScoringResultMessage) error {
		require.Equal(t, expected, received)
		cancel() // Завершаем чтение после первого сообщения
		return nil
	})

	require.True(t, err == nil || ctx.Err() == context.Canceled, "unexpected error: %v", err)
}

func TestKafkaClient_PublishScoringRequest(t *testing.T) {
	require.NoError(t, setupKafkaTopic(requestTopic))

	ctx := context.Background()
	client := NewKafkaClient([]string{kafkaBroker}, requestTopic, resultTopic, testGroupID)

	testMsg := ScoringRequestMessage{
		ApplicationID: "test-app-2",
		UserID:        456,
		Repos:         10,
		Stars:         20,
		Forks:         5,
		Commits:       100,
	}

	err := client.PublishScoringRequest(ctx, testMsg)
	require.NoError(t, err)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   requestTopic,
		GroupID: "test-verification-group",
	})
	defer reader.Close()

	readCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := reader.ReadMessage(readCtx)
	require.NoError(t, err)

	var received ScoringRequestMessage
	err = json.Unmarshal(m.Value, &received)
	require.NoError(t, err)
	require.Equal(t, testMsg, received)
}
