package domain

// AuditDetails mirrors org.egov.land.web.models.AuditDetails (epoch millis).
type AuditDetails struct {
	CreatedBy        string `json:"createdBy"`
	LastModifiedBy   string `json:"lastModifiedBy"`
	CreatedTime      *int64 `json:"createdTime"`
	LastModifiedTime *int64 `json:"lastModifiedTime"`
}
