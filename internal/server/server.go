package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"jwt-auth-service/init/config"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/internal/repository/postgres"
	"jwt-auth-service/internal/server/http/routes"
	"jwt-auth-service/pkg/constants"
	"net/http"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

func NewServer(ctx context.Context, cfg *config.Config, log *logrus.Logger) (*HTTPServer, error) {
	db, err := postgres.NewPostgresConnection(ctx, cfg, log)
	if err != nil {
		return nil, err
	}

	engine := setupGin(cfg)
	router := engine.Group(cfg.ApiEntry)
	routes.InitRouterAndComponents(ctx, router, db, cfg).Routes()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.ApiPort),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        engine,
		MaxHeaderBytes: 1 << 20,
	}

	return &HTTPServer{server: server}, nil
}

func (s *HTTPServer) Run(cfg *config.Config) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err.Error(), constants.ServerCategory)
		}
	}()

	logger.InfoF("success to listen and serve on :%d port", constants.ServerCategory, cfg.ApiPort)

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Error(err.Error(), constants.ServerCategory)
	}
}

func setupGin(cfg *config.Config) *gin.Engine {
	var mode = gin.ReleaseMode
	if cfg.ApiDebug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithFormatter(logger.HTTPLogger))

	return engine
}
