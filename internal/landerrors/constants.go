package landerrors

// Constants ported from org.egov.land.util.LandConstants.
const (
	BPAModule            = "BPA"
	CommonMastersModule  = "common-masters"
	OwnershipCategoryKey = "OwnerShipCategory"

	Citizen = "CITIZEN"

	InvalidSearch              = "INVALID SEARCH"
	InvalidAddress             = "INVALID ADDRESS"
	BoundaryError              = "BOUNDARY ERROR"
	BoundaryMDMSDataError      = "BOUNDARY MDMS DATA ERROR"
	InvalidBoundaryData        = "INVALID BOUNDARY DATA"
	OwnerSearchError           = "OWNER SEARCH ERROR"
	InvalidTenant              = "INVALID TENANT"
	UpdateError                = "UPDATE ERROR"
	InvalidOwnerError          = "INVALID ONWER ERROR"
	IllegalArgumentException   = "ILLEGAL ARGUMENT EXCEPTION"
	BPADuplicateDocument       = "BPA_DUPLICATE_DOCUMENT"
	DuplicateMobileNumber      = "DUPLICATE_MOBILENUMBER_EXCEPTION"
	InvalidTenantIDMDMSKey     = "INVALID TENANTID"
	InvalidTenantIDMDMSMessage = "No data found for this tenentID"

	CommonMasterJSONPath = "$.MdmsRes.common-masters"
	ActiveMasterFilter   = "$.[?(@.active==true)].code"
)
