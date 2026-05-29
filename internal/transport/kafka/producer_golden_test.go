package kafka

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/testutil"
)

func TestSaveLandInfoKafkaPayloadContract(t *testing.T) {
	golden := testutil.ReadGolden(t, "kafka_land_info_request.json")

	var want domain.LandInfoRequest
	if err := json.Unmarshal(golden, &want); err != nil {
		t.Fatal(err)
	}

	mock := &testutil.MockProducer{}
	cfg := config.KafkaConfig{
		SaveLandInfoTopic:   "save-landinfo",
		UpdateLandInfoTopic: "update-landinfo",
	}

	if err := SaveLandInfo(context.Background(), mock, cfg, want); err != nil {
		t.Fatal(err)
	}
	if mock.Topic != "save-landinfo" {
		t.Fatalf("topic=%q", mock.Topic)
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(mock.Payload, &got); err != nil {
		t.Fatal(err)
	}
	if _, ok := got["LandInfo"]; !ok {
		t.Fatalf("kafka payload missing LandInfo key: %s", mock.Payload)
	}
	if _, ok := got["RequestInfo"]; !ok {
		t.Fatalf("kafka payload missing RequestInfo key: %s", mock.Payload)
	}

	// Round-trip must preserve persister-critical casing.
	var round domain.LandInfoRequest
	if err := json.Unmarshal(mock.Payload, &round); err != nil {
		t.Fatal(err)
	}
	if round.LandInfo == nil || round.LandInfo.TenantID != want.LandInfo.TenantID {
		t.Fatalf("tenantId mismatch: %+v", round.LandInfo)
	}
}

func TestUpdateLandInfoUsesUpdateTopic(t *testing.T) {
	mock := &testutil.MockProducer{}
	cfg := config.KafkaConfig{
		SaveLandInfoTopic:   "save-landinfo",
		UpdateLandInfoTopic: "update-landinfo",
	}
	req := domain.LandInfoRequest{
		LandInfo: &domain.LandInfo{ID: "id-1", TenantID: "pb.amritsar"},
	}
	if err := UpdateLandInfo(context.Background(), mock, cfg, req); err != nil {
		t.Fatal(err)
	}
	if mock.Topic != "update-landinfo" {
		t.Fatalf("topic=%q", mock.Topic)
	}
}
