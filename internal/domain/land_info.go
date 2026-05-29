package domain

import "encoding/json"

// LandInfo mirrors org.egov.land.web.models.LandInfo.
type LandInfo struct {
	ID                string          `json:"id,omitempty"`
	LandUID           string          `json:"landUId,omitempty"`
	LandUniqueRegNo   string          `json:"landUniqueRegNo,omitempty"`
	TenantID          string          `json:"tenantId,omitempty"`
	Status            Status          `json:"status,omitempty"`
	Address           *Address        `json:"address,omitempty"`
	OwnershipCategory string          `json:"ownershipCategory,omitempty"`
	Owners            []OwnerInfo     `json:"owners,omitempty"`
	Institution       *Institution    `json:"institution,omitempty"`
	Source            Source          `json:"source,omitempty"`
	Channel           Channel         `json:"channel,omitempty"`
	Documents         []Document      `json:"documents,omitempty"`
	Unit              []Unit          `json:"unit,omitempty"`
	AdditionalDetails json.RawMessage `json:"additionalDetails,omitempty"`
	AuditDetails      *AuditDetails   `json:"auditDetails,omitempty"`
}
