package service

import (
	"context"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/mdms"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/pkg/httpclient"
)

// LandUtil ports org.egov.land.util.LandUtil.
type LandUtil struct {
	cfg    *config.Config
	client *httpclient.Client
}

func NewLandUtil(cfg *config.Config, client *httpclient.Client) *LandUtil {
	return &LandUtil{cfg: cfg, client: client}
}

func (u *LandUtil) GetAuditDetails(by string, isCreate bool) *domain.AuditDetails {
	now := time.Now().UnixMilli()
	if isCreate {
		return &domain.AuditDetails{
			CreatedBy:        by,
			LastModifiedBy:   by,
			CreatedTime:      &now,
			LastModifiedTime: &now,
		}
	}
	return &domain.AuditDetails{
		LastModifiedBy:   by,
		LastModifiedTime: &now,
	}
}

func (u *LandUtil) MDMSCall(ctx context.Context, requestInfo *domain.RequestInfo, tenantID string) (map[string]any, error) {
	req := mdms.BuildRequest(requestInfo, tenantID)
	var resp map[string]any
	if err := u.client.PostJSON(ctx, u.cfg.Egov.MDMS.SearchURL(), req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func auditUserUUID(requestInfo *domain.RequestInfo) string {
	if requestInfo == nil || requestInfo.UserInfo == nil {
		return ""
	}
	return requestInfo.UserInfo.UUID
}
