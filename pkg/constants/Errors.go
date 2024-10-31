package constants

import "errors"

var (
	ApiVarsRequiredError     = errors.New("API port and API entry is required")
	PostgresDSNRequiredError = errors.New("postgres DSN is required")

	UserNotFoundError      = errors.New("user not found")
	UserAlreadyExistsError = errors.New("user already exists")

	RefreshTokenInvalidError  = errors.New(`refresh token invalid`)
	RefreshTokenExpiredError  = errors.New("refresh token expired")
	RefreshTokenNotFoundError = errors.New(`refresh token not found`)
)
