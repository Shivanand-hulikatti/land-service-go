package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/app"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	transportHTTP "github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 30 * time.Second

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
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logrus.Infof("received %s, shutting down HTTP server...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Warnf("HTTP server shutdown: %v", err)
		}
	}()

	logrus.Infof("listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("server stopped: %v", err)
	}
	logrus.Info("HTTP server stopped")
}
