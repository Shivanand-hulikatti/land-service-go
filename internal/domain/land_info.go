package domain

import "encoding/json"

// LandInfo mirrors org.egov.land.web.models.LandInfo.
type LandInfo struct {
	ID                string          `json:"id"`
	LandUID           string          `json:"landUId"`
	LandUniqueRegNo   string          `json:"landUniqueRegNo"`
	TenantID          string          `json:"tenantId"`
	Status            Status          `json:"status"`
	Address           *Address        `json:"address"`
	OwnershipCategory string          `json:"ownershipCategory"`
	Owners            []OwnerInfo     `json:"owners"`
	Institution       *Institution    `json:"institution"`
	Source            Source          `json:"source"`
	Channel           Channel         `json:"channel"`
	Documents         []Document      `json:"documents"`
	Unit              []Unit          `json:"unit"`
	AdditionalDetails json.RawMessage `json:"additionalDetails"`
	AuditDetails      *AuditDetails   `json:"auditDetails"`
}
