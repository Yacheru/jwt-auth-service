package service

import (
	"context"
	"database/sql"
	"errors"
	"jwt-auth-service/pkg/utils"
	"time"

	"jwt-auth-service/init/config"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/internal/repository"
	"jwt-auth-service/pkg/constants"
	"jwt-auth-service/pkg/email"
	"jwt-auth-service/pkg/jwt"
)

type JWT struct {
	jwtPostgres  repository.JWTPostgresRepository
	userPostgres repository.UserPostgresRepository
	jwtRedis     repository.JWTRedisRepository
	userRedis    repository.UserRedisRepository

	email    email.Sender
	hasher   utils.Hasher
	tManager jwt.TokenManager

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
	hasher utils.Hasher,
	cfg *config.Config) *JWT {
	return &JWT{
		jwtPostgres:     jwtPostgres,
		userPostgres:    userPostgres,
		jwtRedis:        jwtRedis,
		userRedis:       userRedis,
		email:           email,
		hasher:          hasher,
		tManager:        tManager,
		accessTokenTTL:  time.Duration(cfg.AccessTokenTTL) * time.Minute,
		refreshTokenTTL: time.Duration(cfg.RefreshTokenTTL) * time.Minute,
	}
}

func (j *JWT) SetSession(ctx context.Context, ipAddr, userID string) (*entities.Tokens, error) {
	var (
		tokens = new(entities.Tokens)
		err    error
	)

	tokens.AccessToken, err = j.tManager.NewAccessToken(ipAddr, userID, j.accessTokenTTL)
	if err != nil {
		logger.Error(err.Error(), constants.ServiceCategory)
		return nil, err
	}

	tokens.RefreshToken = j.tManager.NewRefreshToken()

	session := &entities.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    time.Now().Add(j.refreshTokenTTL).Unix(),
	}

	if err := j.jwtPostgres.SetSession(ctx, userID, session); err != nil {
		return nil, err
	}
	if err := j.jwtRedis.SetSession(ctx, userID, session); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (j *JWT) RefreshTokens(ctx context.Context, refreshToken, newUserIp string) (string, error) {
	user, err := j.userPostgres.GetUserByRefresh(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", constants.RefreshTokenNotFoundError
		}
		return "", err
	}

	if time.Now().Unix() > user.ExpiresIn {
		return "", constants.RefreshTokenExpiredError
	}

	//if user.IpAddr != newUserIp {
	//	if err := j.email.SendMail([]string{user.Email}, []byte("Suspicious activity on your account")); err != nil {
	//		logger.Error(err.Error(), constants.SMTPCategory)
	//	}
	//}

	accessToken, err := j.tManager.NewAccessToken(newUserIp, user.UserID, j.accessTokenTTL)
	if err != nil {
		logger.Error(err.Error(), constants.ServiceCategory)
		return "", err
	}

	return accessToken, nil
}
