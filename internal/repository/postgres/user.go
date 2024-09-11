package postgres

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/pkg/constants"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (a *UserPostgres) GetUserById(ctx *gin.Context, userId string) (*entities.User, error) {
	var user = new(entities.User)

	query := `
		SELECT uuid, email, password, ip, refresh_token, expires_in 
		FROM users 
		WHERE uuid = $1
	`
	err := a.db.GetContext(ctx.Request.Context(), user, query, userId)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)

		return nil, err
	}

	return user, nil
}

func (a *UserPostgres) StoreNewUser(ctx *gin.Context, u *entities.User) error {
	query := `INSERT INTO users (uuid, email, ip, password) VALUES ($1, $2, $3, $4) RETURNING uuid`
	_, err := a.db.ExecContext(ctx.Request.Context(), query, u.UserID, u.Email, u.IpAddr, u.Password)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)

		return err
	}

	return nil
}
