package service

import (
	"context"

	"jwt-auth-service/init/config"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/internal/repository"
	"jwt-auth-service/pkg/email"
	"jwt-auth-service/pkg/jwt"
	"jwt-auth-service/pkg/utils"
)

type JWTService interface {
	SetSession(ctx context.Context, ipAddr, userID string) (*entities.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken, newUserIp string) (string, error)
}

type UserService interface {
	RegisterUser(ctx context.Context, u *entities.User) error
	LoginUser(ctx context.Context, u *entities.UserLogin) (*entities.Tokens, error)
}

type Service struct {
	JWTService
	UserService
}

func NewService(
	repo *repository.Repository,
	manager *jwt.Manager,
	email email.Sender,
	cfg *config.Config,
	hasher utils.Hasher) *Service {

	jwtService := NewJWTService(
		repo.JWTPostgresRepository,
		repo.UserPostgresRepository,
		repo.JWTRedisRepository,
		repo.UserRedisRepository,
		manager,
		email,
		hasher,
		cfg)

	return &Service{
		JWTService: jwtService,
		UserService: NewUserService(
			repo.UserPostgresRepository,
			repo.UserRedisRepository,
			jwtService,
			email,
			hasher),
	}
}
