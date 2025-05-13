package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	consumer *kafka.Reader
	producer *kafka.Writer
}

type KafkaClientInterface interface {
	ConsumeScoringRequests(ctx context.Context, handler func(ScoringRequestMessage) error) error
	ProduceScoringResult(ctx context.Context, result ScoringResultMessage) error
}

func NewKafkaClient(brokers []string, requestTopic, resultTopic, groupID string) *KafkaClient {
	return &KafkaClient{
		consumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   requestTopic,
			GroupID: groupID,
		}),
		producer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    resultTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (kc *KafkaClient) ConsumeScoringRequests(ctx context.Context, handler func(ScoringRequestMessage) error) error {
	for {
		m, err := kc.consumer.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var req ScoringRequestMessage
		if err := json.Unmarshal(m.Value, &req); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			continue
		}

		if err := handler(req); err != nil {
			log.Printf("handler error: %v", err)
		}
	}
}

func (kc *KafkaClient) ProduceScoringResult(ctx context.Context, result ScoringResultMessage) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return kc.producer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}
