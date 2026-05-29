package domain

// AuthUserInfo mirrors org.egov.land.web.models.UserInfo (auth token payload variant).
type AuthUserInfo struct {
	TenantID         string       `json:"tenantId,omitempty"`
	UUID             string       `json:"uuid,omitempty"`
	UserName         string       `json:"userName,omitempty"`
	Password         string       `json:"password,omitempty"`
	IDToken          string       `json:"idToken,omitempty"`
	Mobile           string       `json:"mobile,omitempty"`
	Email            string       `json:"email,omitempty"`
	PrimaryRole      []Role       `json:"primaryrole,omitempty"`
	AdditionalRoles  []TenantRole `json:"additionalroles,omitempty"`
}
