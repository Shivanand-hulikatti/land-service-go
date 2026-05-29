package domain

// LandSearchCriteria mirrors org.egov.land.web.models.LandSearchCriteria.
// Query params are bound via form tags for Gin @ModelAttribute parity.
type LandSearchCriteria struct {
	TenantID     string   `form:"tenantId" json:"tenantId"`
	IDs          []string `form:"ids" json:"ids,omitempty"`
	LandUID      string   `form:"landUId" json:"landUId,omitempty"`
	MobileNumber string   `form:"mobileNumber" json:"mobileNumber,omitempty"`
	Offset       *int     `form:"offset" json:"offset,omitempty"`
	Limit        *int     `form:"limit" json:"limit,omitempty"`
	Locality     string   `form:"locality" json:"locality,omitempty"`

	// Populated internally during mobile-number search; not in API contract.
	UserIDs []string `json:"-"`
}

// IsEmpty mirrors LandSearchCriteria.isEmpty().
func (c LandSearchCriteria) IsEmpty() bool {
	return c.TenantID == "" && len(c.IDs) == 0 && c.LandUID == "" &&
		c.MobileNumber == "" && c.Locality == ""
}

// TenantIDOnly mirrors LandSearchCriteria.tenantIdOnly().
func (c LandSearchCriteria) TenantIDOnly() bool {
	return c.TenantID != "" && len(c.IDs) == 0 && c.LandUID == "" &&
		c.MobileNumber == "" && c.Locality == ""
}
