package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/app"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

// SetupRouter registers routes under the configured context path (e.g. /land-services).
func SetupRouter(engine *gin.Engine, deps *app.Dependencies) {
	base := engine.Group(normalizeContextPath(deps.Config.Server.ContextPath))

	base.GET("/health", func(c *gin.Context) {
		healthHandler(c, deps)
	})

	landHandler := NewLandHandler(deps.LandService)
	land := base.Group("/v1/land")
	land.POST("/_create", landHandler.Create)
	land.POST("/_update", landHandler.Update)
	land.POST("/_search", landHandler.Search)
}

func healthHandler(c *gin.Context, deps *app.Dependencies) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "UP"
	if deps.DB == nil {
		dbStatus = "DOWN"
	} else if err := postgres.Ping(ctx, deps.DB); err != nil {
		dbStatus = "DOWN"
	}

	kafkaStatus := "UP"
	if deps.Producer == nil {
		kafkaStatus = "DOWN"
	}

	overall := http.StatusOK
	if dbStatus == "DOWN" || kafkaStatus == "DOWN" {
		overall = http.StatusServiceUnavailable
	}

	c.JSON(overall, gin.H{
		"status":  statusFromComponents(dbStatus, kafkaStatus),
		"service": "land-services-go",
		"components": gin.H{
			"database": dbStatus,
			"kafka":    kafkaStatus,
		},
	})
}

func statusFromComponents(parts ...string) string {
	for _, p := range parts {
		if p == "DOWN" {
			return "DEGRADED"
		}
	}
	return "UP"
}

func normalizeContextPath(path string) string {
	if path == "" || path == "/" {
		return ""
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return strings.TrimSuffix(path, "/")
}
