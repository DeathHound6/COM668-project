package middleware

import (
	"com668-backend/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryHandler() gin.RecoveryFunc {
	return func(ctx *gin.Context, err any) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &utility.ErrorResponseSchema{
			Error: err.(string),
		})
	}
}
