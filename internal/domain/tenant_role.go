package domain

// TenantRole mirrors org.egov.land.web.models.TenantRole.
type TenantRole struct {
	TenantID string `json:"tenantId,omitempty"`
	Roles    []Role `json:"roles,omitempty"`
}
