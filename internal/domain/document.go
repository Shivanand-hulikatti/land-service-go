package domain

import "encoding/json"

// Document mirrors org.egov.land.web.models.Document.
type Document struct {
	ID                string          `json:"id,omitempty"`
	DocumentType      string          `json:"documentType,omitempty"`
	FileStoreID       string          `json:"fileStoreId,omitempty"`
	DocumentUID       string          `json:"documentUid,omitempty"`
	AdditionalDetails json.RawMessage `json:"additionalDetails,omitempty"`
	AuditDetails      *AuditDetails   `json:"auditDetails,omitempty"`
}
