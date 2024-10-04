package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Create a user
// @Description Create a user
// @Tags Users
// @Accept json
// @Produce json
// @Success 201
// @Router /users [post]
func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tx := ctx.MustGet("transaction").(*gorm.DB)
		var body *utility.UserPostRequestBodySchema
		ctx.BindJSON(&body)

		if err := database.CreateUser(tx, body); err != nil {
			ctx.Error(err)
		}
		ctx.Status(http.StatusCreated)
	}
}

// @Summary Login as a user
// @Description Login as a user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200
// @Router /users/login [post]
func LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
