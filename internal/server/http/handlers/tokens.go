package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"jwt-auth-service/internal/entities"
	"jwt-auth-service/pkg/constants"
	"net/http"
)

// RefreshTokens
// @Summary User RefreshTokens
// @Tags tokens
// @Description sign-in user
// @Accept  json
// @Produce  json
// @Param input body entities.RefreshToken true "refresh tokens"
// @Param        guid    query     string  true  "User ID"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/refresh [post]
func (h *Handlers) RefreshTokens(ctx *gin.Context) {
	var refreshToken = new(entities.RefreshToken)
	if err := ctx.ShouldBindJSON(refreshToken); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Provide refresh token")
		return
	}

	AccessToken, err := h.s.JWTService.RefreshTokens(ctx, refreshToken.RefreshToken, ctx.ClientIP())
	if err != nil {
		if errors.Is(err, constants.RefreshTokenNotFoundError) {
			NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "new tokens", AccessToken)
	return
}
