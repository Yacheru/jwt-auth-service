package service

import (
	"github.com/gin-gonic/gin"
	"jwt-auth-service/init/config"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/internal/repository"
	"jwt-auth-service/pkg/email"
	"jwt-auth-service/pkg/jwt"
	"jwt-auth-service/pkg/utils"
)

type JWTService interface {
	SetSession(ctx *gin.Context, ipAddr, userId string) (*entities.Tokens, error)
	RefreshTokens(ctx *gin.Context, refreshToken, userId, newUserIp string) (*entities.Tokens, error)
}

type UserService interface {
	StoreNewUser(ctx *gin.Context, u *entities.User) error
}

type Service struct {
	JWTService
	UserService
}

func NewService(repo *repository.Repository, manager *jwt.Manager, email email.Sender, cfg *config.Config, hasher utils.Hasher) *Service {
	return &Service{
		JWTService: NewJWTService(
			repo.JWTPostgresRepository,
			repo.UserPostgresRepository,
			repo.JWTRedisRepository,
			repo.UserRedisRepository,
			manager,
			email,
			cfg),
		UserService: NewUserService(
			repo.UserPostgresRepository,
			repo.UserRedisRepository,
			email,
			hasher),
	}
}
