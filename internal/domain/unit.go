package domain

import "encoding/json"

// Unit mirrors org.egov.land.web.models.Unit.
type Unit struct {
	ID                string          `json:"id,omitempty"`
	TenantID          string          `json:"tenantId,omitempty"`
	FloorNo           string          `json:"floorNo,omitempty"`
	UnitType          string          `json:"unitType,omitempty"`
	UsageCategory     string          `json:"usageCategory,omitempty"`
	OccupancyType     string          `json:"occupancyType,omitempty"`
	OccupancyDate     *int64          `json:"occupancyDate,omitempty"`
	AdditionalDetails json.RawMessage `json:"additionalDetails,omitempty"`
	AuditDetails      *AuditDetails   `json:"auditDetails,omitempty"`
}
