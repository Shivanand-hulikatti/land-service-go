package domain

import "encoding/json"

// GeoLocation mirrors org.egov.land.web.models.GeoLocation.
type GeoLocation struct {
	ID                string          `json:"id,omitempty"`
	Latitude          *float64        `json:"latitude,omitempty"`
	Longitude         *float64        `json:"longitude,omitempty"`
	AdditionalDetails json.RawMessage `json:"additionalDetails,omitempty"`
}
