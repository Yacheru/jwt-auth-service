package postgres

import (
	"context"
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

func (a *UserPostgres) GetUserByRefresh(ctx context.Context, refreshToken string) (*entities.User, error) {
	var user = new(entities.User)

	query := `
		SELECT uuid, email, nickname, ip, password, refresh_token, expires_in 
		FROM users
		WHERE refresh_token = $1
	`
	err := a.db.GetContext(ctx, user, query, refreshToken)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return nil, err
	}

	return user, nil
}

func (a *UserPostgres) GetUserID(ctx context.Context, email, password string) (string, error) {
	var uuid string

	query := `
		SELECT uuid
		FROM users 
		WHERE email = $1 AND password = $2
	`
	err := a.db.GetContext(ctx, &uuid, query, email, password)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return "", err
	}

	return uuid, nil
}

func (a *UserPostgres) StoreNewUser(ctx context.Context, u *entities.User) error {
	query := `INSERT INTO users (uuid, email, ip, password, nickname) VALUES ($1, $2, $3, $4, $5) RETURNING uuid`
	_, err := a.db.ExecContext(ctx, query, u.UserID, u.Email, u.IpAddr, u.Password, u.Nickname)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return err
	}

	return nil
}
