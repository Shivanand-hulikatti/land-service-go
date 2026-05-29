package domain

// Status mirrors org.egov.land.web.models.Status.
type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

// Source mirrors org.egov.land.web.models.Source.
type Source string

const (
	SourceMunicipalRecords Source = "MUNICIPAL_RECORDS"
	SourceFieldSurvey      Source = "FIELD_SURVEY"
)

// Channel mirrors org.egov.land.web.models.Channel.
type Channel string

const (
	ChannelSystem     Channel = "SYSTEM"
	ChannelCFCCounter Channel = "CFC_COUNTER"
	ChannelCitizen    Channel = "CITIZEN"
	ChannelDataEntry  Channel = "DATA_ENTRY"
	ChannelMigration  Channel = "MIGRATION"
)

// Relationship mirrors org.egov.land.web.models.Relationship.
type Relationship string

const (
	RelationshipFather  Relationship = "FATHER"
	RelationshipHusband Relationship = "HUSBAND"
)

// OccupancyType mirrors org.egov.land.web.models.OccupancyType.
type OccupancyType string

const (
	OccupancyTypeOwner  OccupancyType = "OWNER"
	OccupancyTypeTenant OccupancyType = "TENANT"
)
