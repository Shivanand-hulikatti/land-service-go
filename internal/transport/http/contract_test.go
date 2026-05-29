package http

import (
	"encoding/json"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/testutil"
)

// TestSuccessResponseMatchesGoldenShape verifies API contract keys required by clients/persister consumers.
func TestSuccessResponseMatchesGoldenShape(t *testing.T) {
	golden := testutil.ReadGolden(t, "land_info_response.json")

	var want map[string]json.RawMessage
	if err := json.Unmarshal(golden, &want); err != nil {
		t.Fatal(err)
	}
	if _, ok := want["ResponseInfo"]; !ok {
		t.Fatal("golden missing ResponseInfo")
	}
	if _, ok := want["LandInfo"]; !ok {
		t.Fatal("golden missing LandInfo array key")
	}

	ts := int64(1700000000000)
	resp := domain.LandInfoResponse{
		ResponseInfo: responseInfoPtr(&domain.RequestInfo{
			APIID: "Rainmaker",
			Ver:   ".01",
			Ts:    &ts,
			MsgID: "20170310130900|en_IN",
		}, true),
		LandInfo: []domain.LandInfo{{
			ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			TenantID: "pb.amritsar",
			Status:   domain.StatusActive,
		}},
	}

	out, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]json.RawMessage
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"ResponseInfo", "LandInfo"} {
		if _, ok := got[key]; !ok {
			t.Fatalf("response missing %q", key)
		}
	}
}
