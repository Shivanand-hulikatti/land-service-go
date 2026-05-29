package domain

// Workflow mirrors org.egov.land.web.models.Workflow.
type Workflow struct {
	Action                string     `json:"action,omitempty"`
	Assignes              []string   `json:"assignes,omitempty"`
	Comments              string     `json:"comments,omitempty"`
	VarificationDocuments []Document `json:"varificationDocuments,omitempty"`
}
