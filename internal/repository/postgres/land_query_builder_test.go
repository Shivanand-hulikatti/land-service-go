package postgres

import (
	"strings"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

func TestBuildSearchQueryTenantExact(t *testing.T) {
	qb := NewLandQueryBuilder(config.EgovConfig{
		Pagination: config.PaginationConfig{DefaultOffset: 0, DefaultLimit: 10, MaxLimit: 50},
	})
	criteria := domain.LandSearchCriteria{
		TenantID: "pb.amritsar",
		IDs:      []string{"id-1", "id-2"},
	}

	query, args := qb.BuildSearchQuery(criteria)
	if !strings.Contains(query, "landInfo.tenantid=?") {
		t.Fatalf("expected exact tenant filter, query=%s", query)
	}
	if !strings.Contains(query, "landInfo.id IN (?,?)") {
		t.Fatalf("expected ids IN clause, query=%s", query)
	}
	if len(args) < 4 {
		t.Fatalf("args=%v", args)
	}
}

func TestBuildSearchQueryStateLevelTenant(t *testing.T) {
	qb := NewLandQueryBuilder(config.EgovConfig{
		Pagination: config.PaginationConfig{DefaultOffset: 0, DefaultLimit: 10, MaxLimit: 50},
	})
	query, args := qb.BuildSearchQuery(domain.LandSearchCriteria{TenantID: "pb"})
	if !strings.Contains(query, "landInfo.tenantid like ?") {
		t.Fatalf("query=%s", query)
	}
	if len(args) == 0 || args[0] != "%pb%" {
		t.Fatalf("args=%v", args)
	}
}

func TestBuildSearchQueryPaginationNoLimitOffset(t *testing.T) {
	qb := NewLandQueryBuilder(config.EgovConfig{
		Pagination: config.PaginationConfig{DefaultOffset: 0, DefaultLimit: 10, MaxLimit: 50},
	})
	_, args := qb.BuildSearchQuery(domain.LandSearchCriteria{TenantID: "pb.amritsar"})
	// last two args: offset=0, limit+offset=50
	if len(args) < 2 {
		t.Fatal("expected pagination args")
	}
	if args[len(args)-1] != 50 {
		t.Fatalf("expected max limit window 50, got %v", args[len(args)-1])
	}
}

func TestBuildSearchQueryUnlimited(t *testing.T) {
	limit := -1
	qb := NewLandQueryBuilder(config.EgovConfig{
		Pagination: config.PaginationConfig{DefaultOffset: 0, DefaultLimit: 10, MaxLimit: 50},
	})
	query, args := qb.BuildSearchQuery(domain.LandSearchCriteria{
		TenantID: "pb.amritsar",
		Limit:    &limit,
	})
	if strings.Contains(query, "offset_ > ?") {
		t.Fatalf("expected pagination clause removed, query=%s", query)
	}
	if len(args) != 1 {
		t.Fatalf("expected only tenant arg, args=%v", args)
	}
}
