package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FormatResponseMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status, statusOk := ctx.Get("Status")
		body, bodyOk := ctx.Get("Body")
		if !statusOk {
			status = http.StatusNoContent
		}
		if !bodyOk || status.(int) == http.StatusNoContent {
			ctx.AbortWithStatus(status.(int))
			return
		}
		ctx.AbortWithStatusJSON(status.(int), body)
	}
}
