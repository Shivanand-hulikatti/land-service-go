package domain

// Boundary mirrors org.egov.land.web.models.Boundary.
type Boundary struct {
	Code             string     `json:"code,omitempty"`
	Name             string     `json:"name,omitempty"`
	Label            string     `json:"label,omitempty"`
	Latitude         string     `json:"latitude,omitempty"`
	Longitude        string     `json:"longitude,omitempty"`
	Children         []Boundary `json:"children,omitempty"`
	MaterializedPath string     `json:"materializedPath,omitempty"`
}
