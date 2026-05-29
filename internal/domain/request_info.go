package domain

// RequestInfo mirrors org.egov.common.contract.request.RequestInfo (DIGIT common contract).
type RequestInfo struct {
	APIID     string        `json:"apiId,omitempty"`
	Ver       string        `json:"ver,omitempty"`
	Ts        *int64        `json:"ts,omitempty"`
	Action    string        `json:"action,omitempty"`
	Did       string        `json:"did,omitempty"`
	Key       string        `json:"key,omitempty"`
	MsgID     string        `json:"msgId,omitempty"`
	AuthToken string        `json:"authToken,omitempty"`
	UserInfo  *ContractUser `json:"userInfo,omitempty"`
}

// ContractUser mirrors fields used from egov-common-contract User / UserInfo in land-services.
type ContractUser struct {
	ID              *int64       `json:"id,omitempty"`
	UUID            string       `json:"uuid,omitempty"`
	UserName        string       `json:"userName,omitempty"`
	Name            string       `json:"name,omitempty"`
	Type            string       `json:"type,omitempty"`
	MobileNumber    string       `json:"mobileNumber,omitempty"`
	EmailID         string       `json:"emailId,omitempty"`
	TenantID        string       `json:"tenantId,omitempty"`
	Roles           []Role       `json:"roles,omitempty"`
	PrimaryRole     []Role       `json:"primaryrole,omitempty"`
	AdditionalRoles []TenantRole `json:"additionalroles,omitempty"`
}
