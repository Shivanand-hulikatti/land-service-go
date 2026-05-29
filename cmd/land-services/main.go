package main

import (
	"fmt"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/app"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	transportHTTP "github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.Infof("starting land-services-go on port %s (context path %s)", cfg.Server.Port, cfg.Server.ContextPath)

	deps, err := app.NewDependencies(cfg)
	if err != nil {
		logrus.Fatalf("failed to wire dependencies: %v", err)
	}
	defer func() {
		if err := deps.Close(); err != nil {
			logrus.Warnf("shutdown: %v", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger(), transportHTTP.ErrorMiddleware())

	transportHTTP.SetupRouter(engine, deps)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := engine.Run(addr); err != nil {
		logrus.Fatalf("server stopped: %v", err)
	}
}
