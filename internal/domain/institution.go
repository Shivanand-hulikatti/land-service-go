package domain

import "encoding/json"

// Institution mirrors org.egov.land.web.models.Institution.
type Institution struct {
	ID                     string          `json:"id,omitempty"`
	TenantID               string          `json:"tenantId,omitempty"`
	Type                   string          `json:"type,omitempty"`
	Designation            string          `json:"designation,omitempty"`
	NameOfAuthorizedPerson string          `json:"nameOfAuthorizedPerson,omitempty"`
	AdditionalDetails      json.RawMessage `json:"additionalDetails,omitempty"`
}
