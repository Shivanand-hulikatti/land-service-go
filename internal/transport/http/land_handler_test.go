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

type mockLandAPI struct {
	createFn func(ctx context.Context, req *domain.LandInfoRequest) (*domain.LandInfo, error)
	updateFn func(ctx context.Context, req *domain.LandInfoRequest) (*domain.LandInfo, error)
	searchFn func(ctx context.Context, criteria domain.LandSearchCriteria, requestInfo *domain.RequestInfo) ([]domain.LandInfo, error)
}

func (m *mockLandAPI) Create(ctx context.Context, req *domain.LandInfoRequest) (*domain.LandInfo, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}

func (m *mockLandAPI) Update(ctx context.Context, req *domain.LandInfoRequest) (*domain.LandInfo, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, req)
	}
	return nil, nil
}

func (m *mockLandAPI) Search(ctx context.Context, criteria domain.LandSearchCriteria, requestInfo *domain.RequestInfo) ([]domain.LandInfo, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, criteria, requestInfo)
	}
	return nil, nil
}

func TestCreateReturns200WithLandInfoArray(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &mockLandAPI{
		createFn: func(_ context.Context, req *domain.LandInfoRequest) (*domain.LandInfo, error) {
			return &domain.LandInfo{ID: "land-1", TenantID: req.LandInfo.TenantID}, nil
		},
	}
	h := NewLandHandler(svc)

	body := `{"RequestInfo":{"apiId":"Rainmaker","msgId":"m1"},"LandInfo":{"tenantId":"pb.amritsar"}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_create", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Create(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}

	var resp domain.LandInfoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.LandInfo) != 1 || resp.LandInfo[0].ID != "land-1" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.ResponseInfo == nil || resp.ResponseInfo.Status != "successful" {
		t.Fatalf("expected successful ResponseInfo, got %+v", resp.ResponseInfo)
	}
}

func TestCreateCustomExceptionReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &mockLandAPI{
		createFn: func(context.Context, *domain.LandInfoRequest) (*domain.LandInfo, error) {
			return nil, landerrors.New(landerrors.InvalidTenant, "state level not allowed")
		},
	}
	h := NewLandHandler(svc)

	body := `{"RequestInfo":{"msgId":"m1"},"LandInfo":{"tenantId":"pb"}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_create", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}

	var resp domain.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Code != landerrors.InvalidTenant {
		t.Fatalf("unexpected errors: %+v", resp.Errors)
	}
	if resp.ResponseInfo.Status != "failed" {
		t.Fatalf("expected failed status, got %s", resp.ResponseInfo.Status)
	}
}

func TestSearchBindsQueryAndBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotCriteria domain.LandSearchCriteria
	svc := &mockLandAPI{
		searchFn: func(_ context.Context, criteria domain.LandSearchCriteria, _ *domain.RequestInfo) ([]domain.LandInfo, error) {
			gotCriteria = criteria
			return []domain.LandInfo{{ID: "x", TenantID: criteria.TenantID}}, nil
		},
	}
	h := NewLandHandler(svc)

	body := `{"RequestInfo":{"msgId":"m1","userInfo":{"type":"EMPLOYEE"}}}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/land/_search?tenantId=pb.amritsar&ids=id1", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Search(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if gotCriteria.TenantID != "pb.amritsar" || len(gotCriteria.IDs) != 1 || gotCriteria.IDs[0] != "id1" {
		t.Fatalf("criteria not bound: %+v", gotCriteria)
	}
}
