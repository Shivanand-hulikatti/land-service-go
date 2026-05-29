package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

func TestGetLandInfoData_Integration(t *testing.T) {
	if os.Getenv("LAND_DB_INTEGRATION") == "" {
		t.Skip("set LAND_DB_INTEGRATION=1 with Postgres running and schema migrated")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}

	db, err := Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseDB(db)

	repo := NewLandRepository(db, nil, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lands, err := repo.GetLandInfoData(ctx, domain.LandSearchCriteria{
		TenantID: "pb",
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = lands
}
