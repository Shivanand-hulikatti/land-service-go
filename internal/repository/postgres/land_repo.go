package postgres

import (
	"context"
	"fmt"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/repository"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/transport/kafka"
	"gorm.io/gorm"
)

var _ repository.LandRepository = (*LandRepository)(nil)

// LandRepository ports org.egov.land.repository.LandRepository.
type LandRepository struct {
	db       *gorm.DB
	producer kafka.Producer
	cfg      *config.Config
	qb       *LandQueryBuilder
}

// NewLandRepository wires search (GORM) and persister publish (Kafka).
func NewLandRepository(db *gorm.DB, producer kafka.Producer, cfg *config.Config) *LandRepository {
	return &LandRepository{
		db:       db,
		producer: producer,
		cfg:      cfg,
		qb:       NewLandQueryBuilder(cfg.Egov),
	}
}

// Save publishes create payload to save-landinfo (no direct DB write).
func (r *LandRepository) Save(ctx context.Context, req domain.LandInfoRequest) error {
	if r.producer == nil {
		return fmt.Errorf("kafka producer is not configured")
	}
	return kafka.SaveLandInfo(ctx, r.producer, r.cfg.Kafka, req)
}

// Update publishes update payload to update-landinfo (no direct DB write).
func (r *LandRepository) Update(ctx context.Context, req domain.LandInfoRequest) error {
	if r.producer == nil {
		return fmt.Errorf("kafka producer is not configured")
	}
	return kafka.UpdateLandInfo(ctx, r.producer, r.cfg.Kafka, req)
}

// GetLandInfoData runs the Java search query via GORM Raw and maps flat rows to LandInfo.
func (r *LandRepository) GetLandInfoData(ctx context.Context, criteria domain.LandSearchCriteria) ([]domain.LandInfo, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not configured")
	}

	rows, err := r.qb.Search(r.db.WithContext(ctx), criteria).Rows()
	if err != nil {
		return nil, fmt.Errorf("land search query: %w", err)
	}

	lands, err := MapLandInfoRows(rows)
	if err != nil {
		return nil, fmt.Errorf("map land rows: %w", err)
	}
	return lands, nil
}
