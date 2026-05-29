package service

import (
	"context"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/pkg/httpclient"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// LandUserService ports org.egov.land.service.LandUserService.
type LandUserService struct {
	cfg    *config.Config
	client *httpclient.Client
}

func NewLandUserService(cfg *config.Config, client *httpclient.Client) *LandUserService {
	return &LandUserService{cfg: cfg, client: client}
}

func (s *LandUserService) ManageUser(ctx context.Context, landRequest *domain.LandInfoRequest) error {
	if landRequest == nil || landRequest.LandInfo == nil {
		return NewCustomException(InvalidTenant, "land info is required")
	}
	landInfo := landRequest.LandInfo
	requestInfo := landRequest.RequestInfo

	for i := range landInfo.Owners {
		owner := &landInfo.Owners[i]
		if owner.MobileNumber == "" {
			return NewCustomException(InvalidOwnerError, "MobileNo is mandatory for ownerInfo")
		}

		if owner.TenantID == "" {
			owner.TenantID = stateTenantID(landInfo.TenantID)
		}

		userDetailResponse, err := s.userExists(ctx, owner, requestInfo)
		if err != nil {
			return err
		}

		if userDetailResponse == nil || len(userDetailResponse.User) == 0 ||
			!owner.CompareWithExistingUser(userDetailResponse.User[0]) {
			role := citizenRole()
			s.addUserDefaultFields(owner.TenantID, role, owner)
			owner.UserName = uuid.New().String()
			owner.OwnerType = Citizen
			userDetailResponse, err = s.userCall(ctx, s.cfg.Egov.User.CreateURL(), domain.CreateUserRequest{
				RequestInfo: requestInfo,
				User:        owner,
			}, userDOBCreateLayout)
			if err != nil {
				return err
			}
			logrus.Debugf("owner created --> %s", userDetailResponse.User[0].UUID)
		}

		if userDetailResponse != nil {
			setOwnerFields(owner, userDetailResponse)
		}
	}
	return nil
}

func (s *LandUserService) GetUser(ctx context.Context, criteria domain.LandSearchCriteria, requestInfo *domain.RequestInfo) (*domain.UserDetailResponse, error) {
	req := domain.UserSearchRequest{
		RequestInfo:  requestInfo,
		TenantID:     stateTenantID(criteria.TenantID),
		MobileNumber: criteria.MobileNumber,
		Active:       boolPtr(true),
		UserType:     Citizen,
	}
	return s.userCall(ctx, s.cfg.Egov.User.SearchURL(), req, userDOBSearchLayout)
}

func (s *LandUserService) GetUsersForLandInfos(ctx context.Context, landInfos []domain.LandInfo) (*domain.UserDetailResponse, error) {
	uuids := make([]string, 0)
	seen := make(map[string]struct{})
	for _, land := range landInfos {
		for _, owner := range land.Owners {
			if owner.UUID == "" {
				continue
			}
			if owner.Status != nil && *owner.Status {
				if _, ok := seen[owner.UUID]; !ok {
					seen[owner.UUID] = struct{}{}
					uuids = append(uuids, owner.UUID)
				}
			}
		}
	}
	req := domain.UserSearchRequest{UUID: uuids}
	return s.userCall(ctx, s.cfg.Egov.User.SearchURL(), req, userDOBSearchLayout)
}

func (s *LandUserService) userExists(ctx context.Context, owner *domain.OwnerInfo, requestInfo *domain.RequestInfo) (*domain.UserDetailResponse, error) {
	req := domain.UserSearchRequest{
		TenantID:     stateTenantID(owner.TenantID),
		MobileNumber: owner.MobileNumber,
	}
	if owner.UUID != "" {
		req.UUID = []string{owner.UUID}
	}
	return s.userCall(ctx, s.cfg.Egov.User.SearchURL(), req, userDOBSearchLayout)
}

func (s *LandUserService) userCall(ctx context.Context, url string, payload any, dobFormat string) (*domain.UserDetailResponse, error) {
	resp, err := s.client.PostJSONMap(ctx, url, payload)
	if err != nil {
		return nil, err
	}
	out, err := decodeUserDetailResponse(resp, dobFormat)
	if err != nil {
		return nil, NewCustomException(IllegalArgumentException, "unable to decode user service response")
	}
	return out, nil
}

func citizenRole() domain.Role {
	return domain.Role{Code: Citizen, Name: "Citizen"}
}

func (s *LandUserService) addUserDefaultFields(tenantID string, role domain.Role, owner *domain.OwnerInfo) {
	owner.Active = boolPtr(true)
	owner.TenantID = tenantID
	owner.Roles = []domain.Role{role}
	owner.Type = Citizen
}

func setOwnerFields(owner *domain.OwnerInfo, userDetailResponse *domain.UserDetailResponse) {
	if len(userDetailResponse.User) == 0 {
		return
	}
	u := userDetailResponse.User[0]
	owner.ID = u.ID
	owner.UUID = u.UUID
	owner.UserName = u.UserName
}
