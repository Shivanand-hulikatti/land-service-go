package testutil

import (
	"context"
	"encoding/json"
)

// MockProducer records the last Kafka publish for tests.
type MockProducer struct {
	Topic   string
	Payload []byte
	Err     error
}

func (m *MockProducer) Push(_ context.Context, topic string, payload any) error {
	if m.Err != nil {
		return m.Err
	}
	m.Topic = topic
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	m.Payload = data
	return nil
}

func (m *MockProducer) Close() error { return nil }
