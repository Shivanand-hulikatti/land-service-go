package postgres

import (
	"context"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

type mockProducer struct {
	topic   string
	payload any
}

func (m *mockProducer) Push(_ context.Context, topic string, payload any) error {
	m.topic = topic
	m.payload = payload
	return nil
}

func (m *mockProducer) Close() error { return nil }

func TestLandRepositorySaveUpdateKafka(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			SaveLandInfoTopic:   "save-landinfo",
			UpdateLandInfoTopic: "update-landinfo",
		},
	}
	producer := &mockProducer{}
	repo := NewLandRepository(nil, producer, cfg)

	req := domain.LandInfoRequest{
		RequestInfo: &domain.RequestInfo{APIID: "test"},
		LandInfo:    &domain.LandInfo{TenantID: "pb.amritsar"},
	}

	if err := repo.Save(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	if producer.topic != "save-landinfo" {
		t.Fatalf("topic=%s", producer.topic)
	}

	if err := repo.Update(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	if producer.topic != "update-landinfo" {
		t.Fatalf("topic=%s", producer.topic)
	}
}

func TestLandRepositorySaveWithoutProducer(t *testing.T) {
	repo := NewLandRepository(nil, nil, &config.Config{})
	err := repo.Save(context.Background(), domain.LandInfoRequest{})
	if err == nil {
		t.Fatal("expected error when producer missing")
	}
}

func TestLandRepositorySearchWithoutDB(t *testing.T) {
	repo := NewLandRepository(nil, nil, &config.Config{})
	_, err := repo.GetLandInfoData(context.Background(), domain.LandSearchCriteria{})
	if err == nil {
		t.Fatal("expected error when db missing")
	}
}
