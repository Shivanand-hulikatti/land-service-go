package domain

import "encoding/json"

// Unit mirrors org.egov.land.web.models.Unit.
type Unit struct {
	ID                string          `json:"id"`
	TenantID          string          `json:"tenantId"`
	FloorNo           string          `json:"floorNo"`
	UnitType          string          `json:"unitType"`
	UsageCategory     string          `json:"usageCategory"`
	OccupancyType     string          `json:"occupancyType"`
	OccupancyDate     *int64          `json:"occupancyDate"`
	AdditionalDetails json.RawMessage `json:"additionalDetails"`
	AuditDetails      *AuditDetails   `json:"auditDetails"`
}
