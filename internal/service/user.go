package service

import (
	"github.com/gin-gonic/gin"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/internal/repository"
	"jwt-auth-service/pkg/email"
	"jwt-auth-service/pkg/utils"
)

type User struct {
	authPostgres repository.UserPostgresRepository
	userRedis    repository.UserRedisRepository
	email        email.Sender
	hasher       utils.Hasher
}

func NewUserService(authPostgres repository.UserPostgresRepository, userRedis repository.UserRedisRepository, email email.Sender, hasher utils.Hasher) *User {
	return &User{authPostgres: authPostgres, userRedis: userRedis, email: email, hasher: hasher}
}

func (a *User) StoreNewUser(ctx *gin.Context, user *entities.User) error {
	user.Password = a.hasher.Hash(user.Password)

	if err := a.authPostgres.StoreNewUser(ctx, user); err != nil {
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
