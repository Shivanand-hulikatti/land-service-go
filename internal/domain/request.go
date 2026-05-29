package domain

// LandInfoRequest mirrors org.egov.land.web.models.LandInfoRequest.
type LandInfoRequest struct {
	RequestInfo *RequestInfo `json:"RequestInfo,omitempty"`
	LandInfo    *LandInfo    `json:"LandInfo,omitempty"`
}

// RequestInfoWrapper mirrors org.egov.land.web.models.RequestInfoWrapper (search body).
type RequestInfoWrapper struct {
	RequestInfo *RequestInfo `json:"RequestInfo,omitempty"`
}

// CreateUserRequest mirrors org.egov.land.web.models.CreateUserRequest.
// Note: uses lowercase requestInfo per Java @JsonProperty.
type CreateUserRequest struct {
	RequestInfo *RequestInfo `json:"requestInfo,omitempty"`
	User        *OwnerInfo   `json:"user,omitempty"`
}
