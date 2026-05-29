package http

import (
	"context"
	"net/http"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/gin-gonic/gin"
)

// landAPI is the handler-facing contract for LandService (enables tests with mocks).
type landAPI interface {
	Create(ctx context.Context, req *domain.LandInfoRequest) (*domain.LandInfo, error)
	Update(ctx context.Context, req *domain.LandInfoRequest) (*domain.LandInfo, error)
	Search(ctx context.Context, criteria domain.LandSearchCriteria, requestInfo *domain.RequestInfo) ([]domain.LandInfo, error)
}

// LandHandler ports org.egov.land.web.controller.LandController.
type LandHandler struct {
	svc landAPI
}

func NewLandHandler(svc landAPI) *LandHandler {
	return &LandHandler{svc: svc}
}

// Create handles POST /v1/land/_create (200 OK on success, per Java).
func (h *LandHandler) Create(c *gin.Context) {
	var req domain.LandInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, err, nil)
		return
	}
	c.Set("requestInfo", req.RequestInfo)

	landInfo, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err, req.RequestInfo)
		return
	}

	c.JSON(http.StatusOK, domain.LandInfoResponse{
		ResponseInfo: responseInfoPtr(req.RequestInfo, true),
		LandInfo:     []domain.LandInfo{*landInfo},
	})
}

// Update handles POST /v1/land/_update.
func (h *LandHandler) Update(c *gin.Context) {
	var req domain.LandInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, err, nil)
		return
	}
	c.Set("requestInfo", req.RequestInfo)

	landInfo, err := h.svc.Update(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err, req.RequestInfo)
		return
	}

	c.JSON(http.StatusOK, domain.LandInfoResponse{
		ResponseInfo: responseInfoPtr(req.RequestInfo, true),
		LandInfo:     []domain.LandInfo{*landInfo},
	})
}

// Search handles POST /v1/land/_search with JSON body (RequestInfo) and query params (criteria).
func (h *LandHandler) Search(c *gin.Context) {
	var body domain.RequestInfoWrapper
	if err := c.ShouldBindJSON(&body); err != nil {
		writeError(c, err, nil)
		return
	}
	c.Set("requestInfo", body.RequestInfo)

	var criteria domain.LandSearchCriteria
	if err := c.ShouldBindQuery(&criteria); err != nil {
		writeError(c, err, body.RequestInfo)
		return
	}

	lands, err := h.svc.Search(c.Request.Context(), criteria, body.RequestInfo)
	if err != nil {
		writeError(c, err, body.RequestInfo)
		return
	}
	if lands == nil {
		lands = []domain.LandInfo{}
	}

	c.JSON(http.StatusOK, domain.LandInfoResponse{
		ResponseInfo: responseInfoPtr(body.RequestInfo, true),
		LandInfo:     lands,
	})
}

func responseInfoPtr(requestInfo *domain.RequestInfo, success bool) *domain.ResponseInfo {
	ri := domain.NewResponseInfoFromRequest(requestInfo, success)
	return &ri
}
