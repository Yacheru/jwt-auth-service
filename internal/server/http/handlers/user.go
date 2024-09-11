package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"jwt-auth-service/internal/entities"
)

// SignUp
// @Summary User SignUp
// @Tags user-auth
// @Description create user accoung
// @Accept  json
// @Produce  json
// @Param input body entities.User true "sign up info"
// @Success 201 {object} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-up [post]
func (h *Handlers) SignUp(ctx *gin.Context) {
	var user = new(entities.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "body is invalid")
		return
	}

	user.IpAddr = ctx.ClientIP()
	user.UserID = uuid.NewString()

	err := h.s.StoreNewUser(ctx, user)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusCreated, "sign-up success", user.UserID)
}

// SignIn
// @Summary User SignIn
// @Tags user-auth
// @Description sign-in user
// @Accept  json
// @Produce  json
// @Param input body entities.UserSignIn true "sign in info"
// @Success 201 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-in [post]
func (h *Handlers) SignIn(ctx *gin.Context) {
	var userSignIn = new(entities.UserSignIn)
	if err := ctx.ShouldBindJSON(userSignIn); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "provide valid user id!")
		return
	}

	userSignIn.IpAddress = ctx.ClientIP()

	tokens, err := h.s.SetSession(ctx, userSignIn.IpAddress, userSignIn.UserID)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, "Failed to create session")
		return
	}

	tokens.RefreshToken = base64.StdEncoding.EncodeToString([]byte(tokens.RefreshToken))

	NewSuccessResponse(ctx, http.StatusOK, "Your tokens", tokens)
}
