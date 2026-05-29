package kafka

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

func TestNewProducerValidation(t *testing.T) {
	_, err := NewProducer(config.KafkaConfig{})
	if err == nil {
		t.Fatal("expected error for empty config")
	}
}

func TestProducerPush_Integration(t *testing.T) {
	brokers := os.Getenv("LAND_KAFKA_INTEGRATION")
	if brokers == "" {
		t.Skip("set LAND_KAFKA_INTEGRATION=localhost:9092 to run kafka integration test")
	}

	cfg := config.KafkaConfig{
		BootstrapServers:    []string{brokers},
		SaveLandInfoTopic:   "save-landinfo",
		UpdateLandInfoTopic: "update-landinfo",
	}

	producer, err := NewProducer(cfg)
	if err != nil {
		t.Fatalf("NewProducer: %v", err)
	}
	defer producer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload := domain.LandInfoRequest{
		RequestInfo: &domain.RequestInfo{APIID: "integration-test"},
		LandInfo:    &domain.LandInfo{TenantID: "pb.amritsar"},
	}

	if err := SaveLandInfo(ctx, producer, cfg, payload); err != nil {
		t.Fatalf("SaveLandInfo: %v", err)
	}
}
