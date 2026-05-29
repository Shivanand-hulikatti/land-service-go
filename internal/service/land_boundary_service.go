package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/pkg/httpclient"
)

// LandBoundaryService ports org.egov.land.service.LandBoundaryService.
type LandBoundaryService struct {
	cfg    *config.Config
	client *httpclient.Client
}

func NewLandBoundaryService(cfg *config.Config, client *httpclient.Client) *LandBoundaryService {
	return &LandBoundaryService{cfg: cfg, client: client}
}

// GetAreaType enriches address.locality from the location boundary search API.
func (s *LandBoundaryService) GetAreaType(ctx context.Context, request *domain.LandInfoRequest, hierarchyTypeCode string) error {
	if request == nil || request.LandInfo == nil {
		return NewCustomException(InvalidAddress, "The address or locality cannot be null")
	}
	if request.LandInfo.Address == nil || request.LandInfo.Address.Locality == nil {
		return NewCustomException(InvalidAddress, "The address or locality cannot be null")
	}

	localityCode := request.LandInfo.Address.Locality.Code
	tenantID := request.LandInfo.TenantID

	endpoint, err := s.boundarySearchURL(tenantID, hierarchyTypeCode, localityCode)
	if err != nil {
		return err
	}

	resp, err := s.client.PostJSONMap(ctx, endpoint, domain.RequestInfoWrapper{
		RequestInfo: request.RequestInfo,
	})
	if err != nil {
		return err
	}
	if len(resp) == 0 {
		return NewCustomException(BoundaryError, "The response from location service is empty or null")
	}

	boundaries := findBoundariesByCode(resp, localityCode)
	if len(boundaries) == 0 {
		return NewCustomException(BoundaryMDMSDataError, "The boundary data was not found")
	}

	boundary := boundaries[0]
	if boundary.Name == "" {
		return NewCustomException(InvalidBoundaryData,
			fmt.Sprintf("The boundary data for the code %s is not available", localityCode))
	}
	request.LandInfo.Address.Locality = &boundary
	return nil
}

func (s *LandBoundaryService) boundarySearchURL(tenantID, hierarchyTypeCode, localityCode string) (string, error) {
	base := s.cfg.Egov.Location.BoundarySearchURL()
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("tenantId", tenantID)
	if hierarchyTypeCode != "" {
		q.Set("hierarchyTypeCode", hierarchyTypeCode)
	}
	q.Set("boundaryType", "Locality")
	q.Set("codes", localityCode)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func findBoundariesByCode(data any, code string) []domain.Boundary {
	var out []domain.Boundary
	switch v := data.(type) {
	case map[string]any:
		if c, _ := v["code"].(string); c == code {
			if b, ok := mapToBoundary(v); ok {
				out = append(out, b)
			}
		}
		if raw, ok := v["boundary"]; ok {
			out = append(out, collectBoundaries(raw, code)...)
		}
		for _, child := range v {
			out = append(out, findBoundariesByCode(child, code)...)
		}
	case []any:
		for _, item := range v {
			out = append(out, findBoundariesByCode(item, code)...)
		}
	}
	return out
}

func collectBoundaries(raw any, code string) []domain.Boundary {
	switch v := raw.(type) {
	case []any:
		var out []domain.Boundary
		for _, item := range v {
			out = append(out, collectBoundaries(item, code)...)
		}
		return out
	case map[string]any:
		if c, _ := v["code"].(string); c == code {
			if b, ok := mapToBoundary(v); ok {
				return []domain.Boundary{b}
			}
		}
	}
	return nil
}

func mapToBoundary(m map[string]any) (domain.Boundary, bool) {
	raw, err := json.Marshal(m)
	if err != nil {
		return domain.Boundary{}, false
	}
	var b domain.Boundary
	if err := json.Unmarshal(raw, &b); err != nil {
		return domain.Boundary{}, false
	}
	return b, b.Code != ""
}
