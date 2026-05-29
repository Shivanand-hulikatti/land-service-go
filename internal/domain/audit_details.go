package domain

// AuditDetails mirrors org.egov.land.web.models.AuditDetails (epoch millis).
type AuditDetails struct {
	CreatedBy        string `json:"createdBy,omitempty"`
	LastModifiedBy   string `json:"lastModifiedBy,omitempty"`
	CreatedTime      *int64 `json:"createdTime,omitempty"`
	LastModifiedTime *int64 `json:"lastModifiedTime,omitempty"`
}
