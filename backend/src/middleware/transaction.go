package middleware

import (
	"com668-backend/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionRequestMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn := database.GetDBConn()
		tx := conn.Begin()
		ctx.Set("transaction", tx)
		ctx.Next()
	}
}

func TransactionResponseMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("transaction").(*gorm.DB)
		errors := ctx.Errors.Errors()
		if len(errors) > 0 {
			tx.Rollback()
		}
		tx.Commit()
		ctx.Next()
	}
}
