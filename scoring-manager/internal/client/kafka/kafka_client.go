package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type ScoringKafkaClient interface {
	PublishScoringRequest(ctx context.Context, msg ScoringRequestMessage) error
	ConsumeScoringResults(ctx context.Context, handle func(ScoringResultMessage) error) error
}

type KafkaClient struct {
	producer *kafka.Writer
	consumer *kafka.Reader
}

func NewKafkaClient(brokers []string, requestTopic, resultTopic, groupID string) *KafkaClient {
	return &KafkaClient{
		producer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    requestTopic,
			Balancer: &kafka.LeastBytes{},
		},
		consumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   resultTopic,
			GroupID: groupID,
		}),
	}
}

func (kc *KafkaClient) PublishScoringRequest(ctx context.Context, msg ScoringRequestMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return kc.producer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ApplicationID),
		Value: data,
		Time:  time.Now(),
	})
}

func (kc *KafkaClient) ConsumeScoringResults(ctx context.Context, handle func(ScoringResultMessage) error) error {
	for {
		m, err := kc.consumer.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var result ScoringResultMessage
		if err = json.Unmarshal(m.Value, &result); err != nil {
			log.Printf("failed to unmarshal scoring result: %v", err)
			continue
		}

		if err = handle(result); err != nil {
			log.Printf("handler error: %v", err)
			continue
		}
	}
}
