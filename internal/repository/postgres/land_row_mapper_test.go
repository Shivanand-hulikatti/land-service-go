package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMapLandInfoRowsSingleLand(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cols := []string{
		"land_id", "landinfo_tenantid", "landinfo_lastmodifiedtime", "landinfo_createdby",
		"landinfo_lastmodifiedby", "landinfo_createdtime", "additionaldetails",
		"latitude", "longitude", "locality", "buildingname", "city", "plotno", "district",
		"region", "state", "country", "landinfo_ad_id", "landmark", "pincode", "doorno", "street",
		"landinfo_geo_loc", "landuid", "land_regno", "status", "ownershipcategory", "source", "channel",
		"landinfo_un_id", "floorno", "unittype", "usagecategory", "occupancytype", "occupancydate",
		"landinfoowner_id", "landinfoowner_uuid", "isprimaryowner", "ownerstatus", "ownershippercentage",
		"institutionid", "relationship", "land_inst_id", "land_inst_type", "designation", "nameOfAuthorizedPerson",
		"landinfo_doc_id", "landinfo_doc_documenttype", "landinfo_doc_filestore", "documentuid",
	}

	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows(cols).AddRow(
			"land-1", "pb.amritsar", int64(2000), "user-1", "user-1", int64(1000), "{}",
			12.5, 77.5, "LOC01", "Tower", "City", "P1", "Dist", "Reg", "State", "IN",
			"addr-1", "LM", "560001", "10", "Main St",
			"geo-1", "uid-1", "reg-1", "ACTIVE", "INDIVIDUAL", "MUNICIPAL_RECORDS", "SYSTEM",
			nil, nil, nil, nil, nil, nil,
			"owner-1", "uuid-1", true, true, 100.0,
			nil, "FATHER",
			nil, nil, nil, nil,
			nil, nil, nil, nil,
		),
	)

	rows, err := db.Query("SELECT 1")
	if err != nil {
		t.Fatal(err)
	}

	lands, err := MapLandInfoRows(rows)
	if err != nil {
		t.Fatal(err)
	}
	if len(lands) != 1 {
		t.Fatalf("len=%d", len(lands))
	}
	if lands[0].ID != "land-1" || lands[0].TenantID != "pb.amritsar" {
		t.Fatalf("land=%+v", lands[0])
	}
	if lands[0].Address == nil || lands[0].Address.City != "City" {
		t.Fatalf("address=%+v", lands[0].Address)
	}
	if len(lands[0].Owners) != 1 || lands[0].Owners[0].OwnerID != "owner-1" {
		t.Fatalf("owners=%+v", lands[0].Owners)
	}
}
