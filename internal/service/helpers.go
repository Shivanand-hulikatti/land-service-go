package service

import (
	"strings"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

func boolPtr(v bool) *bool {
	return &v
}

func stateTenantID(tenantID string) string {
	parts := strings.SplitN(tenantID, ".", 2)
	return parts[0]
}

func isStateLevelTenant(tenantID string) bool {
	return tenantID != "" && !strings.Contains(tenantID, ".")
}

func ownerActive(owner *domain.OwnerInfo) bool {
	return owner != nil && owner.Active != nil && *owner.Active
}
