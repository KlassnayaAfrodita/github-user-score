package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
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
	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+fmt.Sprint(controller.Port))
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
	_ = setupKafkaTopic(requestTopic)
	_ = setupKafkaTopic(resultTopic)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := NewKafkaClient([]string{kafkaBroker}, requestTopic, resultTopic, testGroupID)

	expectedResult := ScoringResultMessage{
		ApplicationID: "test-app",
		UserID:        123,
		Scoring:       99.1,
	}
	resultBytes, _ := json.Marshal(expectedResult)

	go func() {
		writer := kafka.Writer{
			Addr:     kafka.TCP(kafkaBroker),
			Topic:    resultTopic,
			Balancer: &kafka.LeastBytes{},
		}
		defer writer.Close()
		_ = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(expectedResult.ApplicationID),
			Value: resultBytes,
		})
	}()

	err := client.ConsumeScoringResults(ctx, func(result ScoringResultMessage) error {
		if result != expectedResult {
			t.Errorf("received unexpected result: got %+v, want %+v", result, expectedResult)
		}
		cancel() // прекращаем чтение после одного сообщения
		return nil
	})

	if err != nil && ctx.Err() != context.Canceled {
		t.Fatalf("error consuming scoring result: %v", err)
	}
}

func TestKafkaClient_PublishScoringRequest(t *testing.T) {
	_ = setupKafkaTopic(requestTopic)

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
	if err != nil {
		t.Fatalf("failed to publish scoring request: %v", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   requestTopic,
		GroupID: "test-verification-group",
	})
	defer reader.Close()

	readCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := reader.ReadMessage(readCtx)
	if err != nil {
		t.Fatalf("failed to read message from request topic: %v", err)
	}

	var received ScoringRequestMessage
	err = json.Unmarshal(m.Value, &received)
	if err != nil {
		t.Fatalf("failed to unmarshal received message: %v", err)
	}

	if received != testMsg {
		t.Errorf("received message does not match: got %+v, want %+v", received, testMsg)
	}
}
