package domain

import (
	"encoding/json"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/testutil"
)

func TestErrorResponseGolden(t *testing.T) {
	golden := testutil.ReadGolden(t, "error_response.json")

	var want ErrorResponse
	if err := json.Unmarshal(golden, &want); err != nil {
		t.Fatal(err)
	}

	ts := int64(1700000000000)
	got := NewErrorResponse(&RequestInfo{
		APIID: "Rainmaker",
		Ver:   ".01",
		Ts:    &ts,
		MsgID: "20170310130900|en_IN",
	}, []Error{{
		Code:    "INVALID TENANT",
		Message: " Application cannot be create at StateLevel",
	}})

	out, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}

	var round ErrorResponse
	if err := json.Unmarshal(out, &round); err != nil {
		t.Fatal(err)
	}
	if round.ResponseInfo.Status != "failed" || len(round.Errors) != 1 {
		t.Fatalf("got %+v", round)
	}
	if round.Errors[0].Code != want.Errors[0].Code {
		t.Fatalf("code=%q", round.Errors[0].Code)
	}
	if round.ResponseInfo.ResMsgID != "uief87324" {
		t.Fatalf("resMsgId=%q", round.ResponseInfo.ResMsgID)
	}
}
