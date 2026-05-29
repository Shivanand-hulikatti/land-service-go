package domain

import "encoding/json"

// GeoLocation mirrors org.egov.land.web.models.GeoLocation.
type GeoLocation struct {
	ID                string          `json:"id"`
	Latitude          *float64        `json:"latitude"`
	Longitude         *float64        `json:"longitude"`
	AdditionalDetails json.RawMessage `json:"additionalDetails"`
}
