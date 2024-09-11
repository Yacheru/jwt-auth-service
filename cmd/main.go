package main

import (
	"context"
	"jwt-auth-service/init/config"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/internal/server"
	"jwt-auth-service/pkg/constants"
	"os/signal"
	"syscall"
)

// @title Jwt-Auth-Api
// @version 1.0
// @description jwt-auth-service

// @host localhost
// @BasePath /

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cfg := &config.ServerConfig

	if err := config.InitConfig(); err != nil {
		cancel()
	}

	log := logger.InitLogger(cfg.ApiDebug)

	app, err := server.NewServer(ctx, cfg, log)
	if err != nil {
		cancel()
	}
	logger.Info("server configured", constants.MainCategory)

	if app != nil {
		if err := app.Run(cfg); err != nil {
			cancel()
		}
		logger.Info("server is running", constants.MainCategory)
	}

	<-ctx.Done()

	if app != nil {
		app.Shutdown(ctx)

		logger.Info("http-server shutdown", constants.MainCategory)
	}

	logger.Info("service shutdown", constants.MainCategory)
}
