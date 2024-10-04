package middleware

import (
	"com668-backend/utility"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerResponseMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		statusCode := 200
		if len(ctx.Errors) > 0 {
			response := &utility.ErrorResponseSchema{Errors: make([]utility.ErrorSchema, 0)}
			for _, err := range ctx.Errors {
				response.Errors = append(response.Errors, utility.ErrorSchema{
					Message: err.Error(),
				})
			}
			ctx.JSON(statusCode, response)
		}
		ctx.Next()
	}
}
