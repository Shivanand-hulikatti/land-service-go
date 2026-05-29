package domain

// Address mirrors org.egov.land.web.models.Address.
type Address struct {
	TenantID         string        `json:"tenantId,omitempty"`
	DoorNo           string        `json:"doorNo,omitempty"`
	PlotNo           string        `json:"plotNo,omitempty"`
	ID               string        `json:"id,omitempty"`
	Landmark         string        `json:"landmark,omitempty"`
	City             string        `json:"city,omitempty"`
	District         string        `json:"district,omitempty"`
	Region           string        `json:"region,omitempty"`
	State            string        `json:"state,omitempty"`
	Country          string        `json:"country,omitempty"`
	Pincode          string        `json:"pincode,omitempty"`
	AdditionDetails  string        `json:"additionDetails,omitempty"`
	BuildingName     string        `json:"buildingName,omitempty"`
	Street           string        `json:"street,omitempty"`
	Locality         *Boundary     `json:"locality,omitempty"`
	GeoLocation      *GeoLocation  `json:"geoLocation,omitempty"`
	AuditDetails     *AuditDetails `json:"auditDetails,omitempty"`
}
