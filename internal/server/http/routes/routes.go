package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "jwt-auth-service/docs"

	"jwt-auth-service/init/config"
	"jwt-auth-service/internal/repository"
	"jwt-auth-service/internal/repository/redis"
	"jwt-auth-service/internal/server/http/handlers"
	"jwt-auth-service/internal/service"
	"jwt-auth-service/pkg/email"
	"jwt-auth-service/pkg/jwt"
	"jwt-auth-service/pkg/utils"
)

type Router struct {
	router  *gin.RouterGroup
	handler *handlers.Handlers
}

func InitRouterAndComponents(ctx context.Context, router *gin.RouterGroup, db *sqlx.DB, cfg *config.Config) *Router {
	r := redis.NewRedisClient(ctx, cfg)
	repo := repository.NewRepository(db, r, cfg)
	tManager := jwt.NewJWTManager(cfg.Salt)
	hasher := utils.NewSHA512Hasher(cfg.Salt)
	mail := email.NewSMPTServer(cfg.EmailSender, cfg.EmailSenderPassword, cfg.SmtpHost, cfg.SmtpPort)
	serv := service.NewService(repo, tManager, mail, cfg, hasher)
	handler := handlers.NewHandlers(serv, tManager)

	return &Router{
		router:  router,
		handler: handler,
	}
}

func (r *Router) Routes() {
	{
		r.router.POST("register", r.handler.Register)
		r.router.POST("login", r.handler.Login)
		r.router.POST("refresh", r.handler.RefreshTokens)
	}

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
