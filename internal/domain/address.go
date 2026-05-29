package domain

// Address mirrors org.egov.land.web.models.Address.
type Address struct {
	TenantID        string        `json:"tenantId"`
	DoorNo          string        `json:"doorNo"`
	PlotNo          string        `json:"plotNo"`
	ID              string        `json:"id"`
	Landmark        string        `json:"landmark"`
	City            string        `json:"city"`
	District        string        `json:"district"`
	Region          string        `json:"region"`
	State           string        `json:"state"`
	Country         string        `json:"country"`
	Pincode         string        `json:"pincode"`
	AdditionDetails string        `json:"additionDetails"`
	BuildingName    string        `json:"buildingName"`
	Street          string        `json:"street"`
	Locality        *Boundary     `json:"locality"`
	GeoLocation     *GeoLocation  `json:"geoLocation"`
	AuditDetails    *AuditDetails `json:"auditDetails"`
}
