package domain

// Boundary mirrors org.egov.land.web.models.Boundary.
type Boundary struct {
	Code             string     `json:"code"`
	Name             string     `json:"name"`
	Label            string     `json:"label"`
	Latitude         string     `json:"latitude"`
	Longitude        string     `json:"longitude"`
	Children         []Boundary `json:"children"`
	MaterializedPath string     `json:"materializedPath"`
}
