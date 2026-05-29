package app

import (
	"fmt"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/repository"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/repository/postgres"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/service"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/transport/kafka"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/validator"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/pkg/httpclient"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Dependencies holds infrastructure wired in main.
type Dependencies struct {
	Config      *config.Config
	DB          *gorm.DB
	Producer    kafka.Producer
	HTTP        *httpclient.Client
	LandRepo    repository.LandRepository
	LandService *service.LandService
}

// NewDependencies initializes HTTP client, optional DB, and optional Kafka producer.
func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	timeout := time.Duration(cfg.HTTPClient.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	deps := &Dependencies{
		Config: cfg,
		HTTP:   httpclient.New(timeout),
	}

	db, err := postgres.Open(cfg)
	if err != nil {
		logrus.Warnf("database unavailable (search will fail until Postgres is up): %v", err)
	} else {
		deps.DB = db
		logrus.Info("connected to PostgreSQL")
	}

	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		logrus.Warnf("kafka producer unavailable (create/update will fail until Kafka is up): %v", err)
	} else {
		deps.Producer = producer
		logrus.Infof("kafka producer ready (topics: %s, %s)",
			cfg.Kafka.SaveLandInfoTopic, cfg.Kafka.UpdateLandInfoTopic)
	}

	deps.LandRepo = postgres.NewLandRepository(deps.DB, deps.Producer, cfg)

	landUtil := service.NewLandUtil(cfg, deps.HTTP)
	boundarySvc := service.NewLandBoundaryService(cfg, deps.HTTP)
	userSvc := service.NewLandUserService(cfg, deps.HTTP)
	mdmsValidator := validator.NewLandMDMSValidator()
	landValidator := validator.NewLandValidator(mdmsValidator)
	enrichmentSvc := service.NewLandEnrichmentService(landUtil, boundarySvc, cfg, userSvc)
	deps.LandService = service.NewLandService(landValidator, enrichmentSvc, userSvc, deps.LandRepo, landUtil)

	return deps, nil
}

// Close releases database and kafka resources.
func (d *Dependencies) Close() error {
	var errs []error
	if d.DB != nil {
		if err := postgres.CloseDB(d.DB); err != nil {
			errs = append(errs, fmt.Errorf("close db: %w", err))
		}
	}
	if d.Producer != nil {
		if err := d.Producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close kafka: %w", err))
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
