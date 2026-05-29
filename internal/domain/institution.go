package domain

import "encoding/json"

// Institution mirrors org.egov.land.web.models.Institution.
type Institution struct {
	ID                     string          `json:"id"`
	TenantID               string          `json:"tenantId"`
	Type                   string          `json:"type"`
	Designation            string          `json:"designation"`
	NameOfAuthorizedPerson string          `json:"nameOfAuthorizedPerson"`
	AdditionalDetails      json.RawMessage `json:"additionalDetails"`
}
