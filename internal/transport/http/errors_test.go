package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/pkg/httpclient"
)

func TestMapToDIGITErrorsCustomExceptionMap(t *testing.T) {
	err := landerrors.NewMap(map[string]string{
		"INVALID OWNERSHIPCATEGORY": "bad category",
		"MDMS DATA ERROR ":          "missing master",
	})
	status, errs := mapToDIGITErrors(err)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d", status)
	}
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestMapToDIGITErrorsHTTPClient(t *testing.T) {
	err := &httpclient.HTTPError{StatusCode: 502, Body: "upstream", URL: "http://x"}
	status, errs := mapToDIGITErrors(err)
	if status != http.StatusBadGateway {
		t.Fatalf("status=%d", status)
	}
	if errs[0].Code != "EXTERNAL_SERVICE_EXCEPTION" {
		t.Fatalf("code=%s", errs[0].Code)
	}
}

func TestMapToDIGITErrorsGeneric(t *testing.T) {
	status, errs := mapToDIGITErrors(errors.New("boom"))
	if status != http.StatusInternalServerError {
		t.Fatalf("status=%d", status)
	}
	if errs[0].Code != "INTERNAL_SERVER_ERROR" {
		t.Fatalf("code=%s", errs[0].Code)
	}
}
