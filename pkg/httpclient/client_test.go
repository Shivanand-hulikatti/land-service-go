package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostJSONSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method=%s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("content-type=%s", r.Header.Get("Content-Type"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["tenantId"] != "pb.amritsar" {
			t.Fatalf("tenantId=%v", body["tenantId"])
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer server.Close()

	client := New(5 * time.Second)
	var resp map[string]any
	err := client.PostJSON(context.Background(), server.URL, map[string]string{
		"tenantId": "pb.amritsar",
	}, &resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("resp=%v", resp)
	}
}

func TestPostJSONHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"Errors":[{"code":"INVALID"}]}`))
	}))
	defer server.Close()

	client := New(5 * time.Second)
	err := client.PostJSON(context.Background(), server.URL, map[string]string{}, nil)
	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d", httpErr.StatusCode)
	}
}

func TestPostJSONMap(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"user": []any{}})
	}))
	defer server.Close()

	client := New(5 * time.Second)
	out, err := client.PostJSONMap(context.Background(), server.URL, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["user"]; !ok {
		t.Fatalf("out=%v", out)
	}
}
