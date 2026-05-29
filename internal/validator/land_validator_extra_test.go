package validator

import (
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
)

func TestValidateDuplicateDocument(t *testing.T) {
	v := NewLandValidator(NewLandMDMSValidator())
	req := &domain.LandInfoRequest{
		LandInfo: &domain.LandInfo{
			OwnershipCategory: "INDIVIDUAL",
			Documents: []domain.Document{
				{FileStoreID: "fs-1"},
				{FileStoreID: "fs-1"},
			},
			Owners: []domain.OwnerInfo{{MobileNumber: "1"}},
		},
	}
	err := v.ValidateLandInfo(req, validMDMSMap())
	if err == nil {
		t.Fatal("expected duplicate document error")
	}
	ce, ok := err.(*landerrors.CustomException)
	if !ok || ce.Code != landerrors.BPADuplicateDocument {
		t.Fatalf("got %v", err)
	}
}

func validMDMSMap() map[string]any {
	return map[string]any{
		"MdmsRes": map[string]any{
			"common-masters": map[string]any{
				"OwnerShipCategory": []any{map[string]any{"code": "INDIVIDUAL"}},
			},
		},
	}
}
