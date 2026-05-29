package postgres

import (
	"strings"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"gorm.io/gorm"
)

const (
	innerJoin       = " INNER JOIN "
	leftOuterJoin   = " LEFT OUTER JOIN "
	landSearchQuery = "SELECT landInfo.*,landInfoaddress.*,landInfoowner.*,landInfounit.*," +
		"landInfogeolocation.*,landInstitution.*,landInfodoc.*,landInfo.id as land_id,landInfo.tenantid as landInfo_tenantId," +
		"landInfo.lastModifiedTime as landinfo_lastmodifiedtime, landInfo.createdBy as landInfo_createdBy,landInfo.lastModifiedBy as landInfo_lastModifiedBy," +
		"landInfo.createdTime as landInfo_createdTime,landInfo.additionalDetails, " +
		"landInfoaddress.id as landInfo_ad_id,landInfogeolocation.id as landInfo_geo_loc," +
		"landInfoowner.id as landInfoowner_id,landInfoowner.uuid as landInfoowner_uuid,landInfoowner.status as ownerstatus,landInfo.landuniqueregno as land_regno," +
		"landInstitution.type as land_inst_type, landInstitution.id as land_inst_id, " +
		"landInfounit.id as landInfo_un_id, landInfodoc.id as landInfo_doc_id,landInfodoc.documenttype as landInfo_doc_documenttype,landInfodoc.filestoreid as landInfo_doc_filestore" +
		" FROM eg_land_landInfo landInfo" + innerJoin +
		"eg_land_Address landInfoaddress ON landInfoaddress.landInfoId = landInfo.id" + leftOuterJoin +
		"eg_land_institution landInstitution ON landInstitution.landInfoId = landInfo.id" + innerJoin +
		"eg_land_ownerInfo landInfoowner ON landInfoowner.landInfoId = landInfo.id AND landInfoowner.status = true " + leftOuterJoin +
		"eg_land_unit landInfounit ON landInfounit.landInfoId = landInfo.id" + leftOuterJoin +
		"eg_land_document landInfodoc ON landInfodoc.landInfoId = landInfo.id" + leftOuterJoin +
		"eg_land_GeoLocation landInfogeolocation ON landInfogeolocation.addressid = landInfoaddress.id"

	paginationWrapper = "SELECT * FROM " +
		"(SELECT *, DENSE_RANK() OVER (ORDER BY landinfo_lastmodifiedtime DESC) offset_ FROM " +
		"({})" +
		" result) result_offset " +
		"WHERE offset_ > ? AND offset_ <= ?"
)

// LandQueryBuilder ports org.egov.land.repository.querybuilder.LandQueryBuilder.
type LandQueryBuilder struct {
	cfg config.EgovConfig
}

func NewLandQueryBuilder(cfg config.EgovConfig) *LandQueryBuilder {
	return &LandQueryBuilder{cfg: cfg}
}

// Search returns a GORM session that runs the Java-equivalent land search SQL.
func (b *LandQueryBuilder) Search(db *gorm.DB, criteria domain.LandSearchCriteria) *gorm.DB {
	query, args := b.buildSearchQuery(criteria)
	return db.Raw(query, args...)
}

// BuildSearchQuery returns SQL and args (for tests). Placeholders use ? (GORM/pg driver binds them).
func (b *LandQueryBuilder) BuildSearchQuery(criteria domain.LandSearchCriteria) (string, []any) {
	return b.buildSearchQuery(criteria)
}

func (b *LandQueryBuilder) buildSearchQuery(criteria domain.LandSearchCriteria) (string, []any) {
	builder := strings.Builder{}
	builder.WriteString(landSearchQuery)

	args := make([]any, 0, 16)

	if criteria.TenantID != "" {
		if strings.Count(criteria.TenantID, ".") == 0 {
			addClause(len(args), &builder)
			builder.WriteString(" landInfo.tenantid like ?")
			args = append(args, "%"+criteria.TenantID+"%")
		} else {
			addClause(len(args), &builder)
			builder.WriteString(" landInfo.tenantid=? ")
			args = append(args, criteria.TenantID)
		}
	}

	if len(criteria.IDs) > 0 {
		addClause(len(args), &builder)
		builder.WriteString(" landInfo.id IN (")
		builder.WriteString(placeholders(len(criteria.IDs)))
		builder.WriteByte(')')
		for _, id := range criteria.IDs {
			args = append(args, id)
		}
	}

	if len(criteria.UserIDs) > 0 {
		addClause(len(args), &builder)
		builder.WriteString(" landInfoowner.uuid IN (")
		builder.WriteString(placeholders(len(criteria.UserIDs)))
		builder.WriteByte(')')
		for _, id := range criteria.UserIDs {
			args = append(args, id)
		}
	}

	if criteria.LandUID != "" {
		addClause(len(args), &builder)
		builder.WriteString(" landInfo.landuid = ? ")
		args = append(args, criteria.LandUID)
	}

	if criteria.Locality != "" {
		addClause(len(args), &builder)
		builder.WriteString(" landInfoaddress.locality = ? ")
		args = append(args, criteria.Locality)
	}

	return b.wrapPagination(builder.String(), args, criteria)
}

func (b *LandQueryBuilder) wrapPagination(query string, args []any, criteria domain.LandSearchCriteria) (string, []any) {
	limit := b.cfg.Pagination.DefaultLimit
	offset := b.cfg.Pagination.DefaultOffset
	finalQuery := strings.Replace(paginationWrapper, "{}", query, 1)

	if criteria.Limit == nil && criteria.Offset == nil {
		limit = b.cfg.Pagination.MaxLimit
	}

	if criteria.Limit != nil {
		if *criteria.Limit <= b.cfg.Pagination.MaxLimit {
			limit = *criteria.Limit
		} else {
			limit = b.cfg.Pagination.MaxLimit
		}
	}

	if criteria.Offset != nil {
		offset = *criteria.Offset
	}

	if limit == -1 {
		finalQuery = strings.Replace(finalQuery, "WHERE offset_ > ? AND offset_ <= ?", "", 1)
		return finalQuery, args
	}

	args = append(args, offset, limit+offset)
	return finalQuery, args
}

func addClause(existingArgs int, builder *strings.Builder) {
	if existingArgs == 0 {
		builder.WriteString(" WHERE ")
	} else {
		builder.WriteString(" AND")
	}
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	parts := make([]string, n)
	for i := range parts {
		parts[i] = "?"
	}
	return strings.Join(parts, ",")
}
