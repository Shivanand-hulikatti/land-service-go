package postgres

import (
	"database/sql"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

type landAggregate struct {
	land     domain.LandInfo
	ownerIDs map[string]struct{}
	unitIDs  map[string]struct{}
	docIDs   map[string]struct{}
}

// MapLandInfoRows ports org.egov.land.repository.rowmapper.LandRowMapper.
func MapLandInfoRows(rows *sql.Rows) ([]domain.LandInfo, error) {
	defer rows.Close()

	buildingMap := make(map[string]*landAggregate)
	order := make([]string, 0)

	for rows.Next() {
		values, cols, err := scanRowValues(rows)
		if err != nil {
			return nil, err
		}
		idx := columnIndexMap(cols)

		id := valString(values, idx, "land_id")
		if id == "" {
			continue
		}

		agg, ok := buildingMap[id]
		if !ok {
			agg = &landAggregate{
				ownerIDs: make(map[string]struct{}),
				unitIDs:  make(map[string]struct{}),
				docIDs:   make(map[string]struct{}),
			}
			agg.land = buildLandInfo(values, idx, id)
			buildingMap[id] = agg
			order = append(order, id)
		}
		addChildren(values, idx, agg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]domain.LandInfo, 0, len(order))
	for _, id := range order {
		out = append(out, buildingMap[id].land)
	}
	return out, nil
}

func buildLandInfo(row []interface{}, idx map[string]int, id string) domain.LandInfo {
	tenantID := valString(row, idx, "landinfo_tenantid", "landInfo_tenantId")

	createdTime, ctOK := valInt64(row, idx, "landinfo_createdtime", "landInfo_createdTime")
	lastModTime, lmOK := valInt64(row, idx, "landinfo_lastmodifiedtime", "landInfo_lastModifiedTime")

	audit := &domain.AuditDetails{
		CreatedBy:        valString(row, idx, "landinfo_createdby", "landInfo_createdBy"),
		LastModifiedBy:   valString(row, idx, "landinfo_lastmodifiedby", "landInfo_lastModifiedBy"),
		CreatedTime:      int64Ptr(createdTime, ctOK),
		LastModifiedTime: int64Ptr(lastModTime, lmOK),
	}

	lat, latOK := valFloat64(row, idx, "latitude")
	lng, lngOK := valFloat64(row, idx, "longitude")
	var geo *domain.GeoLocation
	if geoID := valString(row, idx, "landinfo_geo_loc", "landInfo_geo_loc"); geoID != "" || latOK || lngOK {
		geo = &domain.GeoLocation{
			ID:        geoID,
			Latitude:  lat,
			Longitude: lng,
		}
	}

	localityCode := valString(row, idx, "locality")
	var locality *domain.Boundary
	if localityCode != "" {
		locality = &domain.Boundary{Code: localityCode}
	}

	address := &domain.Address{
		ID:           valString(row, idx, "landinfo_ad_id", "landInfo_ad_id"),
		TenantID:     tenantID,
		BuildingName: valString(row, idx, "buildingname", "buildingName"),
		City:         valString(row, idx, "city"),
		PlotNo:       valString(row, idx, "plotno", "plotNo"),
		District:     valString(row, idx, "district"),
		Region:       valString(row, idx, "region"),
		State:        valString(row, idx, "state"),
		Country:      valString(row, idx, "country"),
		Landmark:     valString(row, idx, "landmark"),
		Pincode:      valString(row, idx, "pincode"),
		DoorNo:       valString(row, idx, "doorno", "doorNo"),
		Street:       valString(row, idx, "street"),
		Locality:     locality,
		GeoLocation:  geo,
	}

	statusStr := valString(row, idx, "status")
	sourceStr := valString(row, idx, "source")
	channelStr := valString(row, idx, "channel")

	return domain.LandInfo{
		ID:                id,
		LandUID:           valString(row, idx, "landuid", "landUid"),
		LandUniqueRegNo:   valString(row, idx, "land_regno", "landuniqueregno"),
		TenantID:          tenantID,
		Status:            domain.Status(statusStr),
		Address:           address,
		OwnershipCategory: valString(row, idx, "ownershipcategory", "ownershipCategory"),
		Source:            domain.Source(sourceStr),
		Channel:           domain.Channel(channelStr),
		AdditionalDetails: valJSONRaw(row, idx, "additionaldetails", "additionalDetails"),
		AuditDetails:      audit,
		Owners:            []domain.OwnerInfo{},
		Documents:         []domain.Document{},
		Unit:              []domain.Unit{},
	}
}

func addChildren(row []interface{}, idx map[string]int, agg *landAggregate) {
	tenantID := agg.land.TenantID

	createdTime, ctOK := valInt64(row, idx, "landinfo_createdtime", "landInfo_createdTime")
	lastModTime, lmOK := valInt64(row, idx, "landinfo_lastmodifiedtime", "landInfo_lastModifiedTime")
	audit := &domain.AuditDetails{
		CreatedBy:        valString(row, idx, "landinfo_createdby", "landInfo_createdBy"),
		LastModifiedBy:   valString(row, idx, "landinfo_lastmodifiedby", "landInfo_lastModifiedBy"),
		CreatedTime:      int64Ptr(createdTime, ctOK),
		LastModifiedTime: int64Ptr(lastModTime, lmOK),
	}

	if unitID := valString(row, idx, "landinfo_un_id", "landInfo_un_id"); unitID != "" {
		if _, exists := agg.unitIDs[unitID]; !exists {
			occDate, occOK := valInt64(row, idx, "occupancydate", "occupancyDate")
			unit := domain.Unit{
				ID:            unitID,
				TenantID:      tenantID,
				FloorNo:       valString(row, idx, "floorno", "floorNo"),
				UnitType:      valString(row, idx, "unittype", "unitType"),
				UsageCategory: valString(row, idx, "usagecategory", "usageCategory"),
				OccupancyType: valString(row, idx, "occupancytype", "occupancyType"),
				OccupancyDate: int64Ptr(occDate, occOK),
				AuditDetails:  audit,
			}
			agg.land.Unit = append(agg.land.Unit, unit)
			agg.unitIDs[unitID] = struct{}{}
		}
	}

	if ownerID := valString(row, idx, "landinfoowner_id", "landInfoowner_id"); ownerID != "" {
		if _, exists := agg.ownerIDs[ownerID]; !exists {
			isPrimary, _ := valBool(row, idx, "isprimaryowner", "isPrimaryOwner")
			status, _ := valBool(row, idx, "ownerstatus", "status")
			ownershipPct, _ := valFloat64(row, idx, "ownershippercentage", "ownerShipPercentage")
			rel := valString(row, idx, "relationship")
			owner := domain.OwnerInfo{
				TenantID:            tenantID,
				OwnerID:             ownerID,
				UUID:                valString(row, idx, "landinfoowner_uuid", "landInfoowner_uuid"),
				IsPrimaryOwner:      isPrimary,
				OwnerShipPercentage: ownershipPct,
				InstitutionID:       valString(row, idx, "institutionid", "institutionId"),
				AuditDetails:        audit,
				Status:              status,
				Relationship:        domain.Relationship(rel),
			}
			agg.land.Owners = append(agg.land.Owners, owner)
			agg.ownerIDs[ownerID] = struct{}{}
		}
	}

	if instID := valString(row, idx, "land_inst_id"); instID != "" {
		agg.land.Institution = &domain.Institution{
			ID:                     instID,
			Type:                   valString(row, idx, "land_inst_type"),
			TenantID:               tenantID,
			Designation:            valString(row, idx, "designation"),
			NameOfAuthorizedPerson: valString(row, idx, "nameOfAuthorizedPerson", "nameofauthorizedperson"),
		}
	}

	if docID := valString(row, idx, "landinfo_doc_id", "landInfo_doc_id"); docID != "" {
		if _, exists := agg.docIDs[docID]; !exists {
			doc := domain.Document{
				ID:           docID,
				DocumentType: valString(row, idx, "landinfo_doc_documenttype", "landInfo_doc_documenttype"),
				FileStoreID:  valString(row, idx, "landinfo_doc_filestore", "landInfo_doc_filestore"),
				DocumentUID:  valString(row, idx, "documentuid", "documentUid"),
				AuditDetails: audit,
			}
			agg.land.Documents = append(agg.land.Documents, doc)
			agg.docIDs[docID] = struct{}{}
		}
	}
}
