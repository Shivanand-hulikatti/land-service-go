package domain

import "encoding/json"

// OwnerInfo mirrors org.egov.land.web.models.OwnerInfo (land owner + egov-user fields).
type OwnerInfo struct {
	TenantID              string          `json:"tenantId,omitempty"`
	Name                  string          `json:"name,omitempty"`
	OwnerID               string          `json:"ownerId,omitempty"`
	MobileNumber          string          `json:"mobileNumber,omitempty"`
	Gender                string          `json:"gender,omitempty"`
	FatherOrHusbandName   string          `json:"fatherOrHusbandName,omitempty"`
	CorrespondenceAddress string          `json:"correspondenceAddress,omitempty"`
	IsPrimaryOwner        *bool           `json:"isPrimaryOwner,omitempty"`
	OwnerShipPercentage   *float64        `json:"ownerShipPercentage,omitempty"`
	OwnerType             string          `json:"ownerType,omitempty"`
	InstitutionID         string          `json:"institutionId,omitempty"`
	Status                *bool           `json:"status,omitempty"`
	Documents             []Document      `json:"documents,omitempty"`
	Relationship          Relationship    `json:"relationship,omitempty"`
	AdditionalDetails     json.RawMessage `json:"additionalDetails,omitempty"`
	ID                    *int64          `json:"id,omitempty"`
	UUID                  string          `json:"uuid,omitempty"`
	UserName              string          `json:"userName,omitempty"`
	Password              string          `json:"password,omitempty"`
	Salutation            string          `json:"salutation,omitempty"`
	EmailID               string          `json:"emailId,omitempty"`
	AltContactNumber      string          `json:"altContactNumber,omitempty"`
	Pan                   string          `json:"pan,omitempty"`
	AadhaarNumber         string          `json:"aadhaarNumber,omitempty"`
	PermanentAddress      string          `json:"permanentAddress,omitempty"`
	PermanentCity         string          `json:"permanentCity,omitempty"`
	PermanentPinCode      string          `json:"permanentPinCode,omitempty"`
	CorrespondenceCity    string          `json:"correspondenceCity,omitempty"`
	CorrespondencePinCode string          `json:"correspondencePinCode,omitempty"`
	Active                *bool           `json:"active,omitempty"`
	Dob                   *int64          `json:"dob,omitempty"`
	PwdExpiryDate         *int64          `json:"pwdExpiryDate,omitempty"`
	Locale                string          `json:"locale,omitempty"`
	Type                  string          `json:"type,omitempty"`
	Signature             string          `json:"signature,omitempty"`
	AccountLocked         *bool           `json:"accountLocked,omitempty"`
	Roles                 []Role          `json:"roles,omitempty"`
	BloodGroup            string          `json:"bloodGroup,omitempty"`
	IdentificationMark    string          `json:"identificationMark,omitempty"`
	Photo                 string          `json:"photo,omitempty"`
	CreatedBy             string          `json:"createdBy,omitempty"`
	CreatedDate           *int64          `json:"createdDate,omitempty"`
	LastModifiedBy        string          `json:"lastModifiedBy,omitempty"`
	LastModifiedDate      *int64          `json:"lastModifiedDate,omitempty"`
	OtpReference          string          `json:"otpReference,omitempty"`
	AuditDetails          *AuditDetails   `json:"auditDetails,omitempty"`
}
