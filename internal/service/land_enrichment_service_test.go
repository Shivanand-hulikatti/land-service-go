package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/pkg/httpclient"
)

func testConfig(locationURL string) *config.Config {
	return &config.Config{
		Egov: config.EgovConfig{
			Location: config.LocationConfig{
				Host:              locationURL,
				ContextPath:       "",
				Endpoint:          "",
				HierarchyTypeCode: "REVENUE",
			},
			MDMS: config.MDMSConfig{
				Host:           "http://mdms.invalid",
				SearchEndpoint: "/v1/_search",
			},
		},
	}
}

func TestEnrichLandInfoRequest_UpdateSetsDefaultsAndUUIDs(t *testing.T) {
	cfg := testConfig("http://unused")
	util := NewLandUtil(cfg, httpclient.New(time.Second))
	enrichment := NewLandEnrichmentService(util, NewLandBoundaryService(cfg, httpclient.New(time.Second)), cfg, nil)

	req := &domain.LandInfoRequest{
		RequestInfo: &domain.RequestInfo{
			UserInfo: &domain.ContractUser{UUID: "auditor-uuid"},
		},
		LandInfo: &domain.LandInfo{
			ID:       "keep-this-id",
			TenantID: "pb.amritsar",
			Address: &domain.Address{
				GeoLocation: &domain.GeoLocation{},
			},
			Unit:      []domain.Unit{{}},
			Documents: []domain.Document{{}},
			Owners:    []domain.OwnerInfo{{}},
		},
	}

	if err := enrichment.EnrichLandInfoRequest(context.Background(), req, true); err != nil {
		t.Fatal(err)
	}

	land := req.LandInfo
	if land.ID != "keep-this-id" {
		t.Fatalf("update must not replace land id, got %q", land.ID)
	}
	if land.Source != domain.SourceMunicipalRecords || land.Channel != domain.ChannelSystem {
		t.Fatalf("defaults: source=%q channel=%q", land.Source, land.Channel)
	}
	if land.AuditDetails == nil || land.AuditDetails.CreatedBy != "auditor-uuid" {
		t.Fatalf("audit: %+v", land.AuditDetails)
	}
	if land.Address.ID == "" || land.Address.GeoLocation.ID == "" {
		t.Fatal("expected address and geoLocation ids")
	}
	if land.Unit[0].ID == "" || land.Unit[0].TenantID != "pb.amritsar" {
		t.Fatal("expected unit id and tenant")
	}
	if land.Documents[0].ID == "" || land.Owners[0].OwnerID == "" {
		t.Fatal("expected document and owner ids")
	}
}

func TestEnrichLandInfoRequest_CreateEnrichesBoundary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"TenantBoundary": []any{
				map[string]any{
					"boundary": []any{
						map[string]any{
							"code":  "SUN01",
							"name":  "Sunshine Layout",
							"label": "Sunshine Layout",
						},
					},
				},
			},
		})
	}))
	defer srv.Close()

	cfg := testConfig(srv.URL)
	client := httpclient.New(time.Second)
	util := NewLandUtil(cfg, client)
	enrichment := NewLandEnrichmentService(util, NewLandBoundaryService(cfg, client), cfg, nil)

	req := &domain.LandInfoRequest{
		RequestInfo: &domain.RequestInfo{
			UserInfo: &domain.ContractUser{UUID: "auditor-uuid"},
		},
		LandInfo: &domain.LandInfo{
			TenantID: "pb.amritsar",
			Address: &domain.Address{
				Locality: &domain.Boundary{Code: "SUN01"},
			},
		},
	}

	if err := enrichment.EnrichLandInfoRequest(context.Background(), req, false); err != nil {
		t.Fatal(err)
	}
	if req.LandInfo.ID == "" {
		t.Fatal("create must assign land id")
	}
	if req.LandInfo.Address.Locality == nil || req.LandInfo.Address.Locality.Name != "Sunshine Layout" {
		t.Fatalf("locality: %+v", req.LandInfo.Address.Locality)
	}
}

func TestEnrichOwnersMergesUserFields(t *testing.T) {
	lands := []domain.LandInfo{{
		ID: "land-1",
		Owners: []domain.OwnerInfo{{
			UUID: "owner-uuid",
		}},
	}}
	users := &domain.UserDetailResponse{
		User: []domain.OwnerInfo{{
			UUID:         "owner-uuid",
			Name:         "Merged Name",
			MobileNumber: "8888888888",
			EmailID:      "a@b.com",
		}},
	}
	if err := enrichOwners(users, lands); err != nil {
		t.Fatal(err)
	}
	if lands[0].Owners[0].Name != "Merged Name" {
		t.Fatalf("name=%q", lands[0].Owners[0].Name)
	}
}

func TestEnrichOwnersMissingUserReturnsError(t *testing.T) {
	lands := []domain.LandInfo{{
		ID:     "land-1",
		Owners: []domain.OwnerInfo{{UUID: "missing"}},
	}}
	err := enrichOwners(&domain.UserDetailResponse{User: []domain.OwnerInfo{}}, lands)
	if err == nil {
		t.Fatal("expected owner search error")
	}
}

func TestEnrichLandInfoRequest_CreateMissingLocality(t *testing.T) {
	cfg := testConfig("http://unused")
	util := NewLandUtil(cfg, httpclient.New(time.Second))
	enrichment := NewLandEnrichmentService(util, NewLandBoundaryService(cfg, httpclient.New(time.Second)), cfg, nil)

	req := &domain.LandInfoRequest{
		RequestInfo: &domain.RequestInfo{UserInfo: &domain.ContractUser{UUID: "u"}},
		LandInfo:    &domain.LandInfo{TenantID: "pb.amritsar"},
	}

	err := enrichment.EnrichLandInfoRequest(context.Background(), req, false)
	if err == nil {
		t.Fatal("expected invalid address error")
	}
	if ce, ok := err.(*CustomException); !ok || ce.Code != InvalidAddress {
		t.Fatalf("got %v", err)
	}
}
