package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"jwt-auth-service/internal/server/http/handlers"
	"net/http"
)

func ParseQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		guid := ctx.Query("guid")

		if guid == "" {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "Missing guid query")
			return
		}

		_, err := uuid.Parse(guid)
		if err != nil {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "Invalid guid query")

			return
		}

		ctx.Next()
	}
}
