package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/config"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/pkg/logger"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/server"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := initLogger(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}

	logger.Infof("[System] Starting %s %s", cfg.Common.Name, cfg.Common.Version)
	logger.Infof("[System] Config: %+v", cfg)

	router := server.SetupRouter(cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.IP, cfg.Server.Port),
		Handler: router,
	}

	go func() {
		logger.Infof("[System] Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("[System] Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof("[System] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("[System] Server forced to shutdown: %v", err)
	}

	logger.Sync()
	logger.Infof("[System] Server exited")
}

func loadConfig() (*config.Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}
	return config.Load(configPath)
}

func initLogger(cfg *config.Config) error {
	return logger.Init(cfg.Common.LogLevel, "console")
}
