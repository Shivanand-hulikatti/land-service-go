package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
)

// Producer publishes persister messages to Kafka (Java org.egov.land.producer.Producer).
type Producer interface {
	Push(ctx context.Context, topic string, payload any) error
	Close() error
}

type syncProducer struct {
	producer sarama.SyncProducer
	cfg      config.KafkaConfig
}

// NewProducer creates a Sarama sync producer aligned with DIGIT persister requirements.
func NewProducer(cfg config.KafkaConfig) (Producer, error) {
	if len(cfg.BootstrapServers) == 0 {
		return nil, fmt.Errorf("kafka bootstrap servers are required")
	}
	if cfg.SaveLandInfoTopic == "" || cfg.UpdateLandInfoTopic == "" {
		return nil, fmt.Errorf("kafka save/update landinfo topics are required")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	retries := cfg.Producer.Retries
	if retries <= 0 {
		retries = 5
	}
	saramaCfg.Producer.Retry.Max = retries
	// Use a conservative Kafka protocol version for compatibility with
	// local Redpanda setups that may not support newer Metadata API versions.
	saramaCfg.Version = sarama.V1_0_0_0

	producer, err := sarama.NewSyncProducer(cfg.BootstrapServers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &syncProducer{producer: producer, cfg: cfg}, nil
}

// Push serializes payload as JSON and sends to topic (mirrors kafkaTemplate.send).
func (p *syncProducer) Push(ctx context.Context, topic string, payload any) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal kafka payload: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka send to %s: %w", topic, err)
	}
	return nil
}

// SaveLandInfo publishes a create request to save-landinfo.
func SaveLandInfo(ctx context.Context, p Producer, cfg config.KafkaConfig, payload any) error {
	return p.Push(ctx, cfg.SaveLandInfoTopic, payload)
}

// UpdateLandInfo publishes an update request to update-landinfo.
func UpdateLandInfo(ctx context.Context, p Producer, cfg config.KafkaConfig, payload any) error {
	return p.Push(ctx, cfg.UpdateLandInfoTopic, payload)
}

func (p *syncProducer) Close() error {
	return p.producer.Close()
}
