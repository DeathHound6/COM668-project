package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
//
//	@Summary Create a user
//	@Description Create a user
//	@Tags Users
//	@Accept json
//	@Produce json
//	@Param user body utility.UserPostRequestBodySchema true "The request body"
//	@Success 201
//	@Header 201 {string} location "GET URL of the Created User"
//	@Failure 400 {object} utility.ErrorResponseSchema
//	@Failure 403 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /users [post]
func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body *utility.UserPostRequestBodySchema
		if err := ctx.BindJSON(&body); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
		}
		user, err := database.CreateUser(ctx, body)
		if err != nil {
			ctx.AbortWithStatusJSON(ctx.GetInt("errorCode"), &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
		}
		ctx.Header("location", fmt.Sprintf("%s://%s/users/%s", ctx.Request.URL.Scheme, ctx.Request.URL.Host, user.UUID))
		ctx.Status(http.StatusCreated)
	}
}

// LoginUser godoc
//
//	@Summary Login as a user
//	@Description Login as a user
//	@Tags Users
//	@Accept json
//	@Produce json
//	@Success 200
//	@Failure 403 {object} utility.ErrorResponseSchema
//	@Failure 404 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /users/login [post]
func LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
