package constants

import "errors"

var (
	ApiVarsRequiredError     = errors.New("API port and API entry is required")
	PostgresDSNRequiredError = errors.New("postgres DSN is required")

	UserNotFoundError = errors.New("user not found")

	RefreshTokenInvalidError         = errors.New(`refresh token invalid`)
	UserDoesNotHaveRefreshTokenError = errors.New(`user does not have refresh token`)
)
