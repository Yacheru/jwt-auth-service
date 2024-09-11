package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"jwt-auth-service/init/config"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/internal/repository"
	"jwt-auth-service/pkg/constants"
	"jwt-auth-service/pkg/email"
	"jwt-auth-service/pkg/jwt"
	"time"
)

type JWT struct {
	jwtPostgres     repository.JWTPostgresRepository
	userPostgres    repository.UserPostgresRepository
	jwtRedis        repository.JWTRedisRepository
	userRedis       repository.UserRedisRepository
	email           email.Sender
	tManager        jwt.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(
	jwtPostgres repository.JWTPostgresRepository,
	userPostgres repository.UserPostgresRepository,
	jwtRedis repository.JWTRedisRepository,
	userRedis repository.UserRedisRepository,
	tManager jwt.TokenManager,
	email email.Sender,
	cfg *config.Config) *JWT {
	return &JWT{
		jwtPostgres:     jwtPostgres,
		userPostgres:    userPostgres,
		jwtRedis:        jwtRedis,
		userRedis:       userRedis,
		email:           email,
		tManager:        tManager,
		accessTokenTTL:  time.Duration(cfg.AccessTokenTTL) * time.Minute,
		refreshTokenTTL: time.Duration(cfg.RefreshTokenTTL) * time.Minute,
	}
}

func (j *JWT) RefreshTokens(ctx *gin.Context, refreshToken, userId, newUserIp string) (*entities.Tokens, error) {
	user, err := j.userRedis.GetUserById(ctx, userId)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if user == nil || errors.Is(err, redis.Nil) {
		user, err = j.userPostgres.GetUserById(ctx, userId)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.UserNotFoundError
		}

		if err != nil {
			return nil, err
		}

		if user.RefreshToken == "" {
			return nil, constants.UserDoesNotHaveRefreshTokenError
		}

	}

	if err := j.tManager.CheckBcryptMatch(user.RefreshToken, refreshToken); err != nil {
		return nil, constants.RefreshTokenInvalidError
	}

	if user.IpAddr != newUserIp {
		if err := j.email.SendMail([]string{user.Email}, []byte("Suspicious activity on your account")); err != nil {
			logger.Error(err.Error(), constants.SMTPCategory)
		}
	}

	tokens, err := j.SetSession(ctx, user.IpAddr, user.UserID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (j *JWT) SetSession(ctx *gin.Context, ipAddr, userId string) (*entities.Tokens, error) {
	var tokens = new(entities.Tokens)
	var err error

	tokens.AccessToken, err = j.tManager.NewAccessToken(ipAddr, userId, j.accessTokenTTL)
	if err != nil {
		logger.Error(err.Error(), constants.ServiceCategory)

		return nil, err
	}

	tokens.RefreshToken = j.tManager.NewRefreshToken()

	bcrypt, err := j.tManager.GenerateRefreshBcrypt(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}

	session := &entities.Session{
		RefreshToken: bcrypt,
		ExpiresIn:    time.Now().Add(j.refreshTokenTTL).Unix(),
	}

	if err := j.jwtPostgres.SetSession(ctx, userId, session); err != nil {
		return nil, err
	}
	if err := j.jwtRedis.SetSession(ctx, userId, session); err != nil {
		return nil, err
	}

	return tokens, nil
}
