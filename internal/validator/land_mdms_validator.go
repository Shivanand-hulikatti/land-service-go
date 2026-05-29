package validator

import (
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/mdms"
)

// LandMDMSValidator ports org.egov.land.validator.LandMDMSValidator.
type LandMDMSValidator struct{}

func NewLandMDMSValidator() *LandMDMSValidator {
	return &LandMDMSValidator{}
}

func (v *LandMDMSValidator) ValidateMdmsData(req *domain.LandInfoRequest, mdmsData map[string]any) error {
	if req.LandInfo == nil {
		return landerrors.New(landerrors.InvalidTenant, "land info is required")
	}

	masterData, err := mdms.ExtractMasterCodes(mdmsData)
	if err != nil {
		return err
	}

	if err := validateIfMasterPresent([]string{landerrors.OwnershipCategoryKey}, masterData); err != nil {
		return err
	}

	for i := range req.LandInfo.Owners {
		if req.LandInfo.Owners[i].OwnerType == "" {
			req.LandInfo.Owners[i].OwnerType = "NONE"
		}
	}

	errs := make(map[string]string)
	if !contains(masterData[landerrors.OwnershipCategoryKey], req.LandInfo.OwnershipCategory) {
		errs["INVALID OWNERSHIPCATEGORY"] = "The OwnerShipCategory '" + req.LandInfo.OwnershipCategory + "' does not exists"
	}
	if len(errs) > 0 {
		return landerrors.NewMap(errs)
	}
	return nil
}

func validateIfMasterPresent(masterNames []string, codes map[string][]string) error {
	errs := make(map[string]string)
	for _, name := range masterNames {
		if len(codes[name]) == 0 {
			errs["MDMS DATA ERROR "] = "Unable to fetch " + name + " codes from MDMS"
		}
	}
	if len(errs) > 0 {
		return landerrors.NewMap(errs)
	}
	return nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
