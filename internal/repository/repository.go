package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"jwt-auth-service/init/config"
	r "jwt-auth-service/internal/repository/redis"

	"jwt-auth-service/internal/entities"
	"jwt-auth-service/internal/repository/postgres"
)

type JWTPostgresRepository interface {
	SetSession(ctx context.Context, userId string, session *entities.Session) error
}

type UserPostgresRepository interface {
	StoreNewUser(ctx context.Context, u *entities.User) error
	GetUserByRefresh(ctx context.Context, refreshToken string) (*entities.User, error)
	GetUserID(ctx context.Context, email, password string) (string, error)
}

type UserRedisRepository interface {
	StoreNewUser(ctx context.Context, u *entities.User) error
	GetUserById(ctx context.Context, userId string) (*entities.User, error)
}

type JWTRedisRepository interface {
	SetSession(ctx context.Context, userId string, session *entities.Session) error
}

type Repository struct {
	JWTPostgresRepository
	UserPostgresRepository
	UserRedisRepository
	JWTRedisRepository
}

func NewRepository(db *sqlx.DB, redis *redis.Client, cfg *config.Config) *Repository {
	return &Repository{
		JWTPostgresRepository:  postgres.NewJWTPostgres(db),
		UserPostgresRepository: postgres.NewUserPostgres(db),
		UserRedisRepository:    r.NewUserRedis(redis, cfg.RedisTTL),
		JWTRedisRepository:     r.NewJWTRedis(redis),
	}
}
