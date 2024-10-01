package controller

import (
	"com668-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// tx := ctx.MustGet("transaction").(*sql.Tx)
		var user models.User
		ctx.BindJSON(&user)
	}
}

func LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
