package domain

// UserSearchRequest mirrors org.egov.land.web.models.UserSearchRequest.
type UserSearchRequest struct {
	RequestInfo   *RequestInfo `json:"RequestInfo,omitempty"`
	UUID          []string     `json:"uuid,omitempty"`
	ID            []string     `json:"id,omitempty"`
	UserName      string       `json:"userName,omitempty"`
	Name          string       `json:"name,omitempty"`
	MobileNumber  string       `json:"mobileNumber,omitempty"`
	AadhaarNumber string       `json:"aadhaarNumber,omitempty"`
	Pan           string       `json:"pan,omitempty"`
	EmailID       string       `json:"emailId,omitempty"`
	FuzzyLogic    bool         `json:"fuzzyLogic,omitempty"`
	Active        *bool        `json:"active,omitempty"`
	TenantID      string       `json:"tenantId,omitempty"`
	PageSize      int          `json:"pageSize,omitempty"`
	PageNumber    int          `json:"pageNumber,omitempty"`
	Sort          []string     `json:"sort,omitempty"`
	UserType      string       `json:"userType,omitempty"`
	RoleCodes     []string     `json:"roleCodes,omitempty"`
}
