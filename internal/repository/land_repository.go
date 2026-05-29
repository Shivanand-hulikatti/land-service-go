package repository

import (
	"context"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

// LandRepository is the persistence contract used by the service layer.
type LandRepository interface {
	Save(ctx context.Context, req domain.LandInfoRequest) error
	Update(ctx context.Context, req domain.LandInfoRequest) error
	GetLandInfoData(ctx context.Context, criteria domain.LandSearchCriteria) ([]domain.LandInfo, error)
}
