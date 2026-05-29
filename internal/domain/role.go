package domain

// Role mirrors org.egov.land.web.models.Role.
type Role struct {
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
