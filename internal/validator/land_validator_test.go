package validator

import (
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
)

func TestValidateDuplicateUser(t *testing.T) {
	v := NewLandValidator(NewLandMDMSValidator())
	req := &domain.LandInfoRequest{
		LandInfo: &domain.LandInfo{
			OwnershipCategory: "INDIVIDUAL",
			Owners: []domain.OwnerInfo{
				{MobileNumber: "9999999999"},
				{MobileNumber: "9999999999"},
			},
		},
	}
	err := v.ValidateLandInfo(req, map[string]any{
		"MdmsRes": map[string]any{
			"common-masters": map[string]any{
				"OwnerShipCategory": []any{map[string]any{"code": "INDIVIDUAL"}},
			},
		},
	})
	if err == nil {
		t.Fatal("expected duplicate mobile error")
	}
	if ce, ok := err.(*landerrors.CustomException); !ok || ce.Code != landerrors.DuplicateMobileNumber {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSearchCitizenRequiresTenantWhenFiltered(t *testing.T) {
	v := NewLandValidator(NewLandMDMSValidator())
	err := v.ValidateSearch(&domain.RequestInfo{
		UserInfo: &domain.ContractUser{Type: landerrors.Citizen},
	}, domain.LandSearchCriteria{IDs: []string{"id-1"}})
	if err == nil {
		t.Fatal("expected tenant required error")
	}
}

func TestValidateSearchEmployeeEmptyNotAllowed(t *testing.T) {
	v := NewLandValidator(NewLandMDMSValidator())
	err := v.ValidateSearch(&domain.RequestInfo{
		UserInfo: &domain.ContractUser{Type: "EMPLOYEE"},
	}, domain.LandSearchCriteria{})
	if err == nil {
		t.Fatal("expected invalid search error")
	}
}
