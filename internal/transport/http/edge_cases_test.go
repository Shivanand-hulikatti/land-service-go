package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
	"github.com/gin-gonic/gin"
)

// Edge-case matrix coverage (Section 12.5) at the HTTP boundary with a mocked service.

func TestEdgeSearchEmptyLandInfoArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewLandHandler(&mockLandAPI{
		searchFn: func(context.Context, domain.LandSearchCriteria, *domain.RequestInfo) ([]domain.LandInfo, error) {
			return []domain.LandInfo{}, nil
		},
	})

	body := `{"RequestInfo":{"msgId":"m1","userInfo":{"type":"EMPLOYEE"}}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_search?tenantId=pb.amritsar", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Search(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var resp domain.LandInfoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.LandInfo) != 0 {
		t.Fatalf("expected [], got %v", resp.LandInfo)
	}
}

func TestEdgeUpdateMissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewLandHandler(&mockLandAPI{
		updateFn: func(context.Context, *domain.LandInfoRequest) (*domain.LandInfo, error) {
			return nil, landerrors.New(landerrors.UpdateError, "Id is mandatory to update ")
		},
	})

	body := `{"RequestInfo":{"msgId":"m1"},"LandInfo":{"tenantId":"pb.amritsar"}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_update", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Update(c)

	assertDIGITErrorCode(t, w, landerrors.UpdateError)
}

func TestEdgeCreateDuplicateMobile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewLandHandler(&mockLandAPI{
		createFn: func(context.Context, *domain.LandInfoRequest) (*domain.LandInfo, error) {
			return nil, landerrors.New(landerrors.DuplicateMobileNumber, "Duplicate mobile numbers found for owners")
		},
	})

	body := `{"RequestInfo":{"msgId":"m1"},"LandInfo":{"tenantId":"pb.amritsar"}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_create", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Create(c)
	assertDIGITErrorCode(t, w, landerrors.DuplicateMobileNumber)
}

func TestEdgeSearchEmployeeNoParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewLandHandler(&mockLandAPI{
		searchFn: func(context.Context, domain.LandSearchCriteria, *domain.RequestInfo) ([]domain.LandInfo, error) {
			return nil, landerrors.New(landerrors.InvalidSearch, "Search without any paramters is not allowed")
		},
	})

	body := `{"RequestInfo":{"msgId":"m1","userInfo":{"type":"EMPLOYEE"}}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_search", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Search(c)
	assertDIGITErrorCode(t, w, landerrors.InvalidSearch)
}

func TestEdgeInvalidJSONBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewLandHandler(&mockLandAPI{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_create", bytes.NewBufferString("{not-json"))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Create(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", w.Code)
	}
	var resp domain.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ResponseInfo.Status != "failed" {
		t.Fatalf("status=%s", resp.ResponseInfo.Status)
	}
}

func assertDIGITErrorCode(t *testing.T, w *httptest.ResponseRecorder, code string) {
	t.Helper()
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var resp domain.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Errors) == 0 || resp.Errors[0].Code != code {
		t.Fatalf("errors=%+v want code %q", resp.Errors, code)
	}
	if resp.ResponseInfo.Status != "failed" {
		t.Fatalf("ResponseInfo.status=%s", resp.ResponseInfo.Status)
	}
}
