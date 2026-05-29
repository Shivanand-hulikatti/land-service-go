package domain

import "encoding/json"

// ConstructionDetail mirrors org.egov.land.web.models.ConstructionDetail.
type ConstructionDetail struct {
	ID                string          `json:"id,omitempty"`
	CarpetArea        *float64        `json:"carpetArea,omitempty"`
	BuiltUpArea       *float64        `json:"builtUpArea,omitempty"`
	PlinthArea        *float64        `json:"plinthArea,omitempty"`
	SuperBuiltUpArea  *float64        `json:"superBuiltUpArea,omitempty"`
	ConstructionType  string          `json:"constructionType,omitempty"`
	ConstructionDate  *int64          `json:"constructionDate,omitempty"`
	Dimensions        json.RawMessage `json:"dimensions,omitempty"`
	AuditDetails      *AuditDetails   `json:"auditDetails,omitempty"`
	AdditionalDetails json.RawMessage `json:"additionalDetails,omitempty"`
}
