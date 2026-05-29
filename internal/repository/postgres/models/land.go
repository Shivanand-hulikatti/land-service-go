package models

import (
	"encoding/json"
)

// LandInfo maps eg_land_landInfo (GORM read model for search joins).
type LandInfo struct {
	ID                string          `gorm:"column:id;primaryKey"`
	LandUID           string          `gorm:"column:landUid"`
	LandUniqueRegNo   string          `gorm:"column:landUniqueRegNo"`
	TenantID          string          `gorm:"column:tenantId"`
	Status            string          `gorm:"column:status"`
	OwnershipCategory string          `gorm:"column:ownershipCategory"`
	Source            string          `gorm:"column:source"`
	Channel           string          `gorm:"column:channel"`
	AdditionalDetails json.RawMessage `gorm:"column:additionalDetails;type:jsonb"`
	CreatedBy         string          `gorm:"column:createdby"`
	LastModifiedBy    string          `gorm:"column:lastmodifiedby"`
	CreatedTime       *int64          `gorm:"column:createdtime"`
	LastModifiedTime  *int64          `gorm:"column:lastmodifiedtime"`

	Address      *Address      `gorm:"foreignKey:LandInfoID;references:ID"`
	Institution  *Institution  `gorm:"foreignKey:LandInfoID;references:ID"`
	Owners       []OwnerInfo   `gorm:"foreignKey:LandInfoID;references:ID"`
	Units        []Unit        `gorm:"foreignKey:LandInfoID;references:ID"`
	Documents    []Document    `gorm:"foreignKey:LandInfoID;references:ID"`
}

func (LandInfo) TableName() string { return "eg_land_landInfo" }

// Address maps eg_land_Address.
type Address struct {
	ID              string `gorm:"column:id;primaryKey"`
	TenantID        string `gorm:"column:tenantId"`
	DoorNo          string `gorm:"column:doorNo"`
	PlotNo          string `gorm:"column:plotNo"`
	Landmark        string `gorm:"column:landmark"`
	City            string `gorm:"column:city"`
	District        string `gorm:"column:district"`
	Region          string `gorm:"column:region"`
	State           string `gorm:"column:state"`
	Country         string `gorm:"column:country"`
	Locality        string `gorm:"column:locality"`
	Pincode         string `gorm:"column:pincode"`
	BuildingName    string `gorm:"column:buildingName"`
	Street          string `gorm:"column:street"`
	LandInfoID      string `gorm:"column:landInfoId"`
	CreatedBy       string `gorm:"column:createdby"`
	LastModifiedBy  string `gorm:"column:lastmodifiedby"`
	CreatedTime     *int64 `gorm:"column:createdtime"`
	LastModifiedTime *int64 `gorm:"column:lastmodifiedtime"`

	GeoLocation *GeoLocation `gorm:"foreignKey:AddressID;references:ID"`
}

func (Address) TableName() string { return "eg_land_Address" }

// GeoLocation maps eg_land_GeoLocation.
type GeoLocation struct {
	ID               string          `gorm:"column:id;primaryKey"`
	Latitude         *float64        `gorm:"column:latitude"`
	Longitude        *float64        `gorm:"column:longitude"`
	AddressID        string          `gorm:"column:addressId"`
	AdditionalDetails json.RawMessage `gorm:"column:additionalDetails;type:jsonb"`
	CreatedBy        string          `gorm:"column:createdby"`
	LastModifiedBy   string          `gorm:"column:lastmodifiedby"`
	CreatedTime      *int64          `gorm:"column:createdtime"`
	LastModifiedTime *int64          `gorm:"column:lastmodifiedtime"`
}

func (GeoLocation) TableName() string { return "eg_land_GeoLocation" }

// OwnerInfo maps eg_land_ownerInfo.
type OwnerInfo struct {
	ID                  string          `gorm:"column:id;primaryKey"`
	UUID                string          `gorm:"column:uuid"`
	IsPrimaryOwner      *bool           `gorm:"column:isprimaryowner"`
	OwnershipPercentage *float64        `gorm:"column:ownershippercentage"`
	InstitutionID       string          `gorm:"column:institutionId"`
	AdditionalDetails   json.RawMessage `gorm:"column:additionalDetails;type:jsonb"`
	LandInfoID          string          `gorm:"column:landInfoId"`
	Relationship        string          `gorm:"column:relationship"`
	Status              *bool           `gorm:"column:status"`
	CreatedBy           string          `gorm:"column:createdby"`
	LastModifiedBy      string          `gorm:"column:lastmodifiedby"`
	CreatedTime         *int64          `gorm:"column:createdtime"`
	LastModifiedTime    *int64          `gorm:"column:lastmodifiedtime"`
}

func (OwnerInfo) TableName() string { return "eg_land_ownerInfo" }

// Institution maps eg_land_institution.
type Institution struct {
	ID                     string          `gorm:"column:id;primaryKey"`
	TenantID               string          `gorm:"column:tenantId"`
	Type                   string          `gorm:"column:type"`
	Designation            string          `gorm:"column:designation"`
	NameOfAuthorizedPerson string          `gorm:"column:nameOfAuthorizedPerson"`
	AdditionalDetails      json.RawMessage `gorm:"column:additionalDetails;type:jsonb"`
	LandInfoID             string          `gorm:"column:landInfoId"`
	CreatedBy              string          `gorm:"column:createdby"`
	LastModifiedBy         string          `gorm:"column:lastmodifiedby"`
	CreatedTime            *int64          `gorm:"column:createdtime"`
	LastModifiedTime       *int64          `gorm:"column:lastmodifiedtime"`
}

func (Institution) TableName() string { return "eg_land_institution" }

// Document maps eg_land_document.
type Document struct {
	ID               string          `gorm:"column:id;primaryKey"`
	DocumentType     string          `gorm:"column:documentType"`
	FileStoreID      string          `gorm:"column:fileStoreId"`
	DocumentUID      string          `gorm:"column:documentUid"`
	AdditionalDetails json.RawMessage `gorm:"column:additionalDetails;type:jsonb"`
	LandInfoID       string          `gorm:"column:landInfoId"`
	CreatedBy        string          `gorm:"column:createdby"`
	LastModifiedBy   string          `gorm:"column:lastmodifiedby"`
	CreatedTime      *int64          `gorm:"column:createdtime"`
	LastModifiedTime *int64          `gorm:"column:lastmodifiedtime"`
}

func (Document) TableName() string { return "eg_land_document" }

// Unit maps eg_land_unit.
type Unit struct {
	ID               string `gorm:"column:id;primaryKey"`
	TenantID         string `gorm:"column:tenantId"`
	FloorNo          string `gorm:"column:floorNo"`
	UnitType         string `gorm:"column:unitType"`
	UsageCategory    string `gorm:"column:usageCategory"`
	OccupancyType    string `gorm:"column:occupancyType"`
	OccupancyDate    *int64 `gorm:"column:occupancyDate"`
	LandInfoID       string `gorm:"column:landInfoId"`
	CreatedBy        string `gorm:"column:createdby"`
	LastModifiedBy   string `gorm:"column:lastmodifiedby"`
	CreatedTime      *int64 `gorm:"column:createdtime"`
	LastModifiedTime *int64 `gorm:"column:lastmodifiedtime"`
}

func (Unit) TableName() string { return "eg_land_unit" }
