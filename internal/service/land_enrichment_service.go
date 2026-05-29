package service

import (
	"context"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// LandEnrichmentService ports org.egov.land.service.LandEnrichmentService.
type LandEnrichmentService struct {
	util     *LandUtil
	boundary *LandBoundaryService
	cfg      *config.Config
	users    *LandUserService
}

func NewLandEnrichmentService(
	util *LandUtil,
	boundary *LandBoundaryService,
	cfg *config.Config,
	users *LandUserService,
) *LandEnrichmentService {
	return &LandEnrichmentService{
		util:     util,
		boundary: boundary,
		cfg:      cfg,
		users:    users,
	}
}

func (s *LandEnrichmentService) EnrichLandInfoRequest(ctx context.Context, landRequest *domain.LandInfoRequest, isUpdate bool) error {
	if landRequest == nil || landRequest.LandInfo == nil || landRequest.RequestInfo == nil {
		return NewCustomException(InvalidTenant, "land info is required")
	}

	auditDetails := s.util.GetAuditDetails(auditUserUUID(landRequest.RequestInfo), true)
	landRequest.LandInfo.AuditDetails = auditDetails

	if err := s.enrichForCreate(ctx, landRequest, isUpdate); err != nil {
		return err
	}

	enrichLandInfoChildren(landRequest.LandInfo, auditDetails)
	return nil
}

func (s *LandEnrichmentService) enrichForCreate(ctx context.Context, landRequest *domain.LandInfoRequest, isUpdate bool) error {
	if isUpdate {
		return nil
	}
	landRequest.LandInfo.ID = uuid.New().String()
	return s.boundary.GetAreaType(ctx, landRequest, s.cfg.Egov.Location.HierarchyTypeCode)
}

func enrichLandInfoChildren(land *domain.LandInfo, auditDetails *domain.AuditDetails) {
	enrichInstitution(land)
	applyLandDefaults(land)
	enrichAddress(land, auditDetails)
	enrichUnits(land, auditDetails)
	enrichDocuments(land, auditDetails)
	assignOwnerIDsAndAudit(land, auditDetails)
}

func enrichInstitution(land *domain.LandInfo) {
	if land.Institution == nil {
		return
	}
	if land.Institution.ID == "" {
		land.Institution.ID = uuid.New().String()
	}
	if land.Institution.TenantID == "" {
		land.Institution.TenantID = land.TenantID
	}
}

func applyLandDefaults(land *domain.LandInfo) {
	if land.Channel == "" {
		land.Channel = domain.ChannelSystem
	}
	if land.Source == "" {
		land.Source = domain.SourceMunicipalRecords
	}
}

func enrichAddress(land *domain.LandInfo, auditDetails *domain.AuditDetails) {
	if land.Address == nil {
		return
	}
	if land.Address.ID == "" {
		land.Address.ID = uuid.New().String()
	}
	land.Address.TenantID = land.TenantID
	land.Address.AuditDetails = auditDetails
	if land.Address.GeoLocation != nil && land.Address.GeoLocation.ID == "" {
		land.Address.GeoLocation.ID = uuid.New().String()
	}
}

func enrichUnits(land *domain.LandInfo, auditDetails *domain.AuditDetails) {
	for i := range land.Unit {
		if land.Unit[i].ID == "" {
			land.Unit[i].ID = uuid.New().String()
		}
		land.Unit[i].TenantID = land.TenantID
		land.Unit[i].AuditDetails = auditDetails
	}
}

func enrichDocuments(land *domain.LandInfo, auditDetails *domain.AuditDetails) {
	for i := range land.Documents {
		if land.Documents[i].ID == "" {
			land.Documents[i].ID = uuid.New().String()
		}
		land.Documents[i].AuditDetails = auditDetails
	}
}

func assignOwnerIDsAndAudit(land *domain.LandInfo, auditDetails *domain.AuditDetails) {
	for i := range land.Owners {
		if land.Owners[i].OwnerID == "" {
			land.Owners[i].OwnerID = uuid.New().String()
		}
		land.Owners[i].AuditDetails = auditDetails
	}
}

func (s *LandEnrichmentService) EnrichLandInfoSearch(
	ctx context.Context,
	landInfos []domain.LandInfo,
	criteria domain.LandSearchCriteria,
	requestInfo *domain.RequestInfo,
) ([]domain.LandInfo, error) {
	if len(landInfos) == 0 {
		return landInfos, nil
	}

	requests := make([]domain.LandInfoRequest, len(landInfos))
	for i := range landInfos {
		requests[i] = domain.LandInfoRequest{RequestInfo: requestInfo, LandInfo: &landInfos[i]}
	}

	if criteria.Limit == nil || *criteria.Limit != -1 {
		for i := range requests {
			if err := s.boundary.GetAreaType(ctx, &requests[i], s.cfg.Egov.Location.HierarchyTypeCode); err != nil {
				return nil, err
			}
		}
	}

	userDetailResponse, err := s.users.GetUsersForLandInfos(ctx, landInfos)
	if err != nil {
		return nil, err
	}
	if err := enrichOwners(userDetailResponse, landInfos); err != nil {
		return nil, err
	}

	if len(landInfos) > 0 && len(landInfos[0].Owners) > 0 {
		logrus.Debug("In enrich service...... ")
	}
	return landInfos, nil
}

func enrichOwners(userDetailResponse *domain.UserDetailResponse, landInfos []domain.LandInfo) error {
	userByUUID := make(map[string]domain.OwnerInfo, len(userDetailResponse.User))
	for _, user := range userDetailResponse.User {
		userByUUID[user.UUID] = user
	}

	for i := range landInfos {
		for j := range landInfos[i].Owners {
			owner := &landInfos[i].Owners[j]
			existing, ok := userByUUID[owner.UUID]
			if !ok {
				return NewCustomException(OwnerSearchError,
					"The owner of the landInfo "+landInfos[i].ID+" is not coming in user search")
			}
			owner.MergeUserWithoutAuditDetail(existing)
		}
	}
	return nil
}
