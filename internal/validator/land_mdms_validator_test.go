package validator

import (
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
)

func TestValidateMdmsDataInvalidOwnershipCategory(t *testing.T) {
	v := NewLandMDMSValidator()
	req := &domain.LandInfoRequest{
		LandInfo: &domain.LandInfo{
			OwnershipCategory: "NOT_A_REAL_CATEGORY",
			Owners:            []domain.OwnerInfo{{}},
		},
	}
	mdms := map[string]any{
		"MdmsRes": map[string]any{
			"common-masters": map[string]any{
				"OwnerShipCategory": []any{map[string]any{"code": "INDIVIDUAL"}},
			},
		},
	}
	err := v.ValidateMdmsData(req, mdms)
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*landerrors.CustomException)
	if !ok || len(ce.Errors) == 0 {
		t.Fatalf("got %v", err)
	}
}

func TestValidateMdmsDataSetsDefaultOwnerType(t *testing.T) {
	v := NewLandMDMSValidator()
	req := &domain.LandInfoRequest{
		LandInfo: &domain.LandInfo{
			OwnershipCategory: "INDIVIDUAL",
			Owners:            []domain.OwnerInfo{{}},
		},
	}
	mdms := map[string]any{
		"MdmsRes": map[string]any{
			"common-masters": map[string]any{
				"OwnerShipCategory": []any{map[string]any{"code": "INDIVIDUAL"}},
			},
		},
	}
	if err := v.ValidateMdmsData(req, mdms); err != nil {
		t.Fatal(err)
	}
	if req.LandInfo.Owners[0].OwnerType != "NONE" {
		t.Fatalf("ownerType=%q", req.LandInfo.Owners[0].OwnerType)
	}
}
