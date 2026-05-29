package service

import (
	"context"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/repository"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/validator"
	"github.com/sirupsen/logrus"
)

type mdmsCaller interface {
	MDMSCall(ctx context.Context, requestInfo *domain.RequestInfo, tenantID string) (map[string]any, error)
}

// LandService ports org.egov.land.service.LandService.
type LandService struct {
	validator  *validator.LandValidator
	enrichment *LandEnrichmentService
	users      *LandUserService
	repo       repository.LandRepository
	mdms       mdmsCaller
}

func NewLandService(
	landValidator *validator.LandValidator,
	enrichment *LandEnrichmentService,
	users *LandUserService,
	repo repository.LandRepository,
	mdms mdmsCaller,
) *LandService {
	return &LandService{
		validator:  landValidator,
		enrichment: enrichment,
		users:      users,
		repo:       repo,
		mdms:       mdms,
	}
}

func (s *LandService) Create(ctx context.Context, landRequest *domain.LandInfoRequest) (*domain.LandInfo, error) {
	if landRequest == nil || landRequest.LandInfo == nil {
		return nil, NewCustomException(InvalidTenant, "land info is required")
	}

	mdmsData, err := s.mdms.MDMSCall(ctx, landRequest.RequestInfo, landRequest.LandInfo.TenantID)
	if err != nil {
		return nil, err
	}
	if isStateLevelTenant(landRequest.LandInfo.TenantID) {
		return nil, NewCustomException(InvalidTenant, " Application cannot be create at StateLevel")
	}

	if err := s.validator.ValidateLandInfo(landRequest, mdmsData); err != nil {
		return nil, err
	}
	if err := s.users.ManageUser(ctx, landRequest); err != nil {
		return nil, err
	}
	if err := s.enrichment.EnrichLandInfoRequest(ctx, landRequest, false); err != nil {
		return nil, err
	}

	setOwnerStatusFromActive(landRequest.LandInfo.Owners)

	if err := s.repo.Save(ctx, *landRequest); err != nil {
		return nil, err
	}
	return landRequest.LandInfo, nil
}

func (s *LandService) Update(ctx context.Context, landRequest *domain.LandInfoRequest) (*domain.LandInfo, error) {
	if landRequest == nil || landRequest.LandInfo == nil {
		return nil, NewCustomException(InvalidTenant, "land info is required")
	}
	landInfo := landRequest.LandInfo

	mdmsData, err := s.mdms.MDMSCall(ctx, landRequest.RequestInfo, landInfo.TenantID)
	if err != nil {
		return nil, err
	}
	if landInfo.ID == "" {
		return nil, NewCustomException(UpdateError, "Id is mandatory to update ")
	}

	defaultOwnerTypes(landInfo.Owners)

	if err := s.validator.ValidateLandInfo(landRequest, mdmsData); err != nil {
		return nil, err
	}
	if err := s.users.ManageUser(ctx, landRequest); err != nil {
		return nil, err
	}
	if err := s.enrichment.EnrichLandInfoRequest(ctx, landRequest, true); err != nil {
		return nil, err
	}

	setOwnerStatusFromActive(landInfo.Owners)

	if err := s.repo.Update(ctx, *landRequest); err != nil {
		return nil, err
	}

	landInfo.Owners = filterActiveOwners(landInfo.Owners)
	return landInfo, nil
}

func defaultOwnerTypes(owners []domain.OwnerInfo) {
	for i := range owners {
		if owners[i].OwnerType == "" {
			owners[i].OwnerType = "NONE"
		}
	}
}

func filterActiveOwners(owners []domain.OwnerInfo) []domain.OwnerInfo {
	if len(owners) <= 1 {
		return owners
	}
	active := make([]domain.OwnerInfo, 0, len(owners))
	for _, owner := range owners {
		if owner.Status != nil && *owner.Status {
			active = append(active, owner)
		}
	}
	return active
}

func (s *LandService) Search(
	ctx context.Context,
	criteria domain.LandSearchCriteria,
	requestInfo *domain.RequestInfo,
) ([]domain.LandInfo, error) {
	if err := s.validator.ValidateSearch(requestInfo, criteria); err != nil {
		return nil, err
	}

	if criteria.MobileNumber != "" {
		lands, err := s.getLandFromMobileNumber(ctx, criteria, requestInfo)
		if err != nil {
			return nil, err
		}
		if len(lands) == 0 {
			return []domain.LandInfo{}, nil
		}
		ids := make([]string, 0, len(lands))
		for _, li := range lands {
			ids = append(ids, li.ID)
		}
		criteria.MobileNumber = ""
		criteria.IDs = ids
	}

	return s.fetchLandInfoData(ctx, criteria, requestInfo)
}

func (s *LandService) getLandFromMobileNumber(
	ctx context.Context,
	criteria domain.LandSearchCriteria,
	requestInfo *domain.RequestInfo,
) ([]domain.LandInfo, error) {
	userDetailResponse, err := s.users.GetUser(ctx, criteria, requestInfo)
	if err != nil {
		return nil, err
	}
	if len(userDetailResponse.User) == 0 {
		return []domain.LandInfo{}, nil
	}

	ids := make([]string, 0, len(userDetailResponse.User))
	for _, u := range userDetailResponse.User {
		ids = append(ids, u.UUID)
	}
	criteria.UserIDs = ids

	landInfo, err := s.repo.GetLandInfoData(ctx, criteria)
	if err != nil {
		return nil, err
	}
	if len(landInfo) == 0 {
		return []domain.LandInfo{}, nil
	}
	return s.enrichment.EnrichLandInfoSearch(ctx, landInfo, criteria, requestInfo)
}

func (s *LandService) fetchLandInfoData(
	ctx context.Context,
	criteria domain.LandSearchCriteria,
	requestInfo *domain.RequestInfo,
) ([]domain.LandInfo, error) {
	landInfos, err := s.repo.GetLandInfoData(ctx, criteria)
	if err != nil {
		return nil, err
	}
	if len(landInfos) == 0 {
		return []domain.LandInfo{}, nil
	}

	logrus.Debug("Received final landInfo response..")
	landInfos, err = s.enrichment.EnrichLandInfoSearch(ctx, landInfos, criteria, requestInfo)
	if err != nil {
		return nil, err
	}
	if len(landInfos) > 0 {
		logrus.Debug("Received final landInfo response after enrichment..")
	}
	return landInfos, nil
}

func setOwnerStatusFromActive(owners []domain.OwnerInfo) {
	for i := range owners {
		if ownerActive(&owners[i]) {
			owners[i].Status = boolPtr(true)
		} else {
			owners[i].Status = boolPtr(false)
		}
	}
}
