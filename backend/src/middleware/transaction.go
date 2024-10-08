package middleware

import (
	"com668-backend/database"
	"com668-backend/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TransactionRequestMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn := database.GetDBConn()
		tx := conn.Begin()
		if tx.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, &utility.ErrorResponseSchema{
				Error: tx.Error.Error(),
			})
		}
		ctx.Set("transaction", tx)
		tx.Set("context", ctx)
		ctx.Next()
	}
}

func TransactionResponseMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := database.GetDBTransaction(ctx)
		if tx.Error != nil {
			tx.Rollback()
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, &utility.ErrorResponseSchema{
				Error: tx.Error.Error(),
			})
		}
		tx.Commit()
		ctx.Next()
	}
}
