package service

import (
	"context"
	"database/sql"
	"errors"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/internal/repository"
	"jwt-auth-service/pkg/constants"
	"jwt-auth-service/pkg/email"
	"jwt-auth-service/pkg/utils"
	"strings"
)

type User struct {
	userPostgres repository.UserPostgresRepository
	userRedis    repository.UserRedisRepository

	jwtService JWTService

	email  email.Sender
	hasher utils.Hasher
}

func NewUserService(
	userPostgres repository.UserPostgresRepository,
	userRedis repository.UserRedisRepository,
	jwtService JWTService,
	email email.Sender,
	hasher utils.Hasher) *User {
	return &User{userPostgres: userPostgres, userRedis: userRedis, jwtService: jwtService, email: email, hasher: hasher}
}

func (a *User) RegisterUser(ctx context.Context, user *entities.User) error {
	user.Password = a.hasher.Hash(user.Password)

	if err := a.userPostgres.StoreNewUser(ctx, user); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return constants.UserAlreadyExistsError
		}
		return err
	}

	if err := a.userRedis.StoreNewUser(ctx, user); err != nil {
		return err
	}

	if err := a.email.SendMail([]string{user.Email}, []byte("You have successfully registered")); err != nil {
		return err
	}

	return nil
}

func (a *User) LoginUser(ctx context.Context, u *entities.UserLogin) (*entities.Tokens, error) {
	u.Password = a.hasher.Hash(u.Password)

	id, err := a.userPostgres.GetUserID(ctx, u.Email, u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.UserNotFoundError
		}
		return nil, err
	}

	tokens, err := a.jwtService.SetSession(ctx, u.IpAddress, id)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
