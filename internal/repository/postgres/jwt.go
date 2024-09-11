package postgres

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/pkg/constants"
)

type JWTPostgres struct {
	db *sqlx.DB
}

func NewJWTPostgres(db *sqlx.DB) *JWTPostgres {
	return &JWTPostgres{db: db}
}

func (j *JWTPostgres) SetSession(ctx *gin.Context, userId string, session *entities.Session) error {
	query := `UPDATE users SET refresh_token = $1, expires_in = $2 WHERE uuid = $3;`
	_, err := j.db.ExecContext(ctx.Request.Context(), query, session.RefreshToken, session.ExpiresIn, userId)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)

		return err
	}

	return nil
}
