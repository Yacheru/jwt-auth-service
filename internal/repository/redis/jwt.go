package redis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/pkg/constants"
)

type JWTRedis struct {
	client *redis.Client
}

func NewJWTRedis(client *redis.Client) *JWTRedis {
	return &JWTRedis{client: client}
}

func (j *JWTRedis) SetSession(ctx *gin.Context, userId string, session *entities.Session) error {
	var user = new(entities.User)

	bytes, err := j.client.Get(ctx, userId).Bytes()
	if err != nil {
		logger.Error(err.Error(), constants.RedisCategory)
		return err
	}

	err = json.Unmarshal(bytes, user)
	if err != nil {
		logger.Error(err.Error(), constants.RedisCategory)
		return err
	}

	user.RefreshToken = session.RefreshToken
	user.ExpiresIn = session.ExpiresIn

	userBytes, err := json.Marshal(user)
	if err != nil {
		logger.Error(err.Error(), constants.RedisCategory)
		return err
	}

	if err := j.client.Set(ctx, userId, userBytes, -1).Err(); err != nil {
		logger.Error(err.Error(), constants.RedisCategory)

		return err
	}

	return nil
}
