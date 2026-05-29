package validator

import (
	"strings"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
)

// LandValidator ports org.egov.land.validator.LandValidator.
type LandValidator struct {
	mdms *LandMDMSValidator
}

func NewLandValidator(mdms *LandMDMSValidator) *LandValidator {
	return &LandValidator{mdms: mdms}
}

func (v *LandValidator) ValidateLandInfo(req *domain.LandInfoRequest, mdmsData map[string]any) error {
	if err := v.mdms.ValidateMdmsData(req, mdmsData); err != nil {
		return err
	}
	if err := validateApplicationDocuments(req); err != nil {
		return err
	}
	return validateDuplicateUser(req)
}

func validateApplicationDocuments(req *domain.LandInfoRequest) error {
	if req.LandInfo == nil || len(req.LandInfo.Documents) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	for _, doc := range req.LandInfo.Documents {
		if _, ok := seen[doc.FileStoreID]; ok {
			return landerrors.New(landerrors.BPADuplicateDocument, "Same document cannot be used multiple times")
		}
		seen[doc.FileStoreID] = struct{}{}
	}
	return nil
}

func validateDuplicateUser(req *domain.LandInfoRequest) error {
	if req.LandInfo == nil || len(req.LandInfo.Owners) <= 1 {
		return nil
	}
	seen := make(map[string]struct{})
	for _, owner := range req.LandInfo.Owners {
		if _, ok := seen[owner.MobileNumber]; ok {
			return landerrors.New(landerrors.DuplicateMobileNumber, "Duplicate mobile numbers found for owners")
		}
		seen[owner.MobileNumber] = struct{}{}
	}
	return nil
}

func (v *LandValidator) ValidateSearch(requestInfo *domain.RequestInfo, criteria domain.LandSearchCriteria) error {
	if requestInfo == nil || requestInfo.UserInfo == nil {
		return landerrors.New(landerrors.InvalidSearch, "RequestInfo user is required")
	}
	userType := requestInfo.UserInfo.Type
	isCitizen := equalIgnoreCase(userType, landerrors.Citizen)

	if !isCitizen && criteria.IsEmpty() {
		return landerrors.New(landerrors.InvalidSearch, "Search without any paramters is not allowed")
	}
	if !isCitizen && !criteria.TenantIDOnly() && criteria.TenantID == "" {
		return landerrors.New(landerrors.InvalidSearch, "TenantId is mandatory in search")
	}
	if isCitizen && !criteria.IsEmpty() && !criteria.TenantIDOnly() && criteria.TenantID == "" {
		return landerrors.New(landerrors.InvalidSearch, "TenantId is mandatory in search")
	}
	return nil
}

func equalIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}
