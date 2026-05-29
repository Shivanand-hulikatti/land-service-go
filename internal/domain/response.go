package domain

// LandInfoResponse mirrors org.egov.land.web.models.LandInfoResponse.
type LandInfoResponse struct {
	ResponseInfo *ResponseInfo `json:"ResponseInfo,omitempty"`
	LandInfo     []LandInfo    `json:"LandInfo"`
}

// UserDetailResponse mirrors org.egov.land.web.models.UserDetailResponse.
type UserDetailResponse struct {
	ResponseInfo *ResponseInfo `json:"responseInfo,omitempty"`
	User         []OwnerInfo   `json:"user,omitempty"`
}
