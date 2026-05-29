package domain

import "encoding/json"

// Document mirrors org.egov.land.web.models.Document.
type Document struct {
	ID                string          `json:"id"`
	DocumentType      string          `json:"documentType"`
	FileStoreID       string          `json:"fileStoreId"`
	DocumentUID       string          `json:"documentUid"`
	AdditionalDetails json.RawMessage `json:"additionalDetails"`
	AuditDetails      *AuditDetails   `json:"auditDetails"`
}
