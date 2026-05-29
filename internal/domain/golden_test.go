package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func goldenPath(name string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "docs", "golden", name)
}

func readGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(goldenPath(name))
	if err != nil {
		t.Fatalf("read golden %s: %v", name, err)
	}
	return data
}

func assertJSONRoundTrip(t *testing.T, input []byte, dest any) {
	t.Helper()
	if err := json.Unmarshal(input, dest); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	out, err := json.Marshal(dest)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	again := reflect.New(reflect.TypeOf(dest).Elem()).Interface()
	if err := json.Unmarshal(out, again); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if !reflect.DeepEqual(dest, again) {
		t.Fatalf("round-trip changed value\ngot:  %+v\nwant: %+v", again, dest)
	}
}

func TestLandInfoRequestGolden(t *testing.T) {
	var req LandInfoRequest
	assertJSONRoundTrip(t, readGolden(t, "land_info_request.json"), &req)

	if req.RequestInfo == nil || req.RequestInfo.APIID != "Rainmaker" {
		t.Fatal("expected RequestInfo.apiId Rainmaker")
	}
	if req.LandInfo == nil || req.LandInfo.TenantID != "pb.amritsar" {
		t.Fatal("expected LandInfo.tenantId pb.amritsar")
	}
	if len(req.LandInfo.Owners) != 1 || req.LandInfo.Owners[0].MobileNumber != "9999999999" {
		t.Fatal("expected one owner with mobileNumber")
	}
}

func TestLandInfoResponseGolden(t *testing.T) {
	var resp LandInfoResponse
	assertJSONRoundTrip(t, readGolden(t, "land_info_response.json"), &resp)

	if resp.ResponseInfo == nil || resp.ResponseInfo.Status != "successful" {
		t.Fatal("expected successful ResponseInfo")
	}
	if len(resp.LandInfo) != 1 || resp.LandInfo[0].Status != StatusActive {
		t.Fatal("expected one ACTIVE LandInfo")
	}
}

func TestRequestInfoWrapperGolden(t *testing.T) {
	var wrap RequestInfoWrapper
	assertJSONRoundTrip(t, readGolden(t, "request_info_wrapper.json"), &wrap)

	if wrap.RequestInfo == nil || wrap.RequestInfo.UserInfo == nil {
		t.Fatal("expected RequestInfo with userInfo")
	}
	if wrap.RequestInfo.UserInfo.Type != "CITIZEN" {
		t.Fatalf("expected CITIZEN user type, got %q", wrap.RequestInfo.UserInfo.Type)
	}
}

func TestNewResponseInfoFromRequest(t *testing.T) {
	ts := int64(1700000000000)
	req := &RequestInfo{
		APIID: "Rainmaker",
		Ver:   ".01",
		Ts:    &ts,
		MsgID: "20170310130900|en_IN",
	}

	resp := NewResponseInfoFromRequest(req, true)
	if resp.Status != "successful" {
		t.Fatalf("status=%q", resp.Status)
	}
	if resp.ResMsgID != "uief87324" {
		t.Fatalf("resMsgId=%q", resp.ResMsgID)
	}
	if resp.APIID != "Rainmaker" || resp.MsgID != "20170310130900|en_IN" {
		t.Fatal("expected apiId and msgId copied from request")
	}

	fail := NewResponseInfoFromRequest(req, false)
	if fail.Status != "failed" {
		t.Fatalf("status=%q", fail.Status)
	}
}

func TestLandSearchCriteriaHelpers(t *testing.T) {
	empty := LandSearchCriteria{}
	if !empty.IsEmpty() {
		t.Fatal("expected empty criteria")
	}

	tenantOnly := LandSearchCriteria{TenantID: "pb.amritsar"}
	if !tenantOnly.TenantIDOnly() {
		t.Fatal("expected tenantIdOnly")
	}
	if tenantOnly.IsEmpty() {
		t.Fatal("expected non-empty when tenantId set")
	}
}
