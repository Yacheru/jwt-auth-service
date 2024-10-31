package handlers

import (
	"errors"
	"jwt-auth-service/pkg/constants"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"jwt-auth-service/internal/entities"
)

// Register
// @Summary User SignUp
// @Tags user-auth
// @Description register account
// @Accept  json
// @Produce  json
// @Param input body entities.User true "sign up info"
// @Success 201 {object} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/register [post]
func (h *Handlers) Register(ctx *gin.Context) {
	var user = new(entities.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "body is invalid")
		return
	}

	user.IpAddr = ctx.ClientIP()
	user.UserID = uuid.NewString()

	err := h.s.UserService.RegisterUser(ctx, user)
	if err != nil {
		if errors.Is(err, constants.UserAlreadyExistsError) {
			NewErrorResponse(ctx, http.StatusConflict, "User already exists")
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusCreated, "register success", user.UserID)
	return
}

// Login
// @Summary User SignIn
// @Tags user-auth
// @Description login user
// @Accept  json
// @Produce  json
// @Param input body entities.UserLogin true "sign in info"
// @Success 201 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/login [post]
func (h *Handlers) Login(ctx *gin.Context) {
	var userLogin = new(entities.UserLogin)
	if err := ctx.ShouldBindJSON(userLogin); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "provide valid user id!")
		return
	}

	userLogin.IpAddress = ctx.ClientIP()

	tokens, err := h.s.UserService.LoginUser(ctx.Request.Context(), userLogin)
	if err != nil {
		if errors.Is(err, constants.UserNotFoundError) {
			NewErrorResponse(ctx, http.StatusUnauthorized, "user not found")
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "client tokens", tokens)
	return
}
