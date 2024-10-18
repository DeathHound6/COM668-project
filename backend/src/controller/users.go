package controller

import (
	"com668-backend/database"
	"com668-backend/middleware"
	"com668-backend/utility"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
//	@Header 201 {string} Location "GET URL of the created User"
//	@Failure 400 {object} utility.ErrorResponseSchema
//	@Failure 403 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /users [post]
func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body *utility.UserPostRequestBodySchema
		if err := ctx.BindJSON(&body); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		user, err := database.CreateUser(ctx, body)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		ctx.Header("Location", fmt.Sprintf("%s://%s/users/%s", ctx.Request.URL.Scheme, ctx.Request.URL.Host, user.UUID))
		ctx.Set("Status", http.StatusCreated)
	}
}

// LoginUser godoc
//
//	@Summary Login as a user
//	@Description Login as a user
//	@Tags Users
//	@Accept json
//	@Produce json
//	@Param request_body body utility.UserLoginRequestBodySchema true "Request Body"
//	@Success 204
//	@Failure 403 {object} utility.ErrorResponseSchema
//	@Failure 404 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /users/login [post]
func LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body *utility.UserLoginRequestBodySchema
		if err := ctx.BindJSON(&body); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		token := jwt.New(middleware.JWTSigningMethod)
		jwtString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
		if err != nil {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		ctx.Set("Status", http.StatusNoContent)
		ctx.Header(middleware.AuthHeaderNameString, fmt.Sprintf("Bearer %s", jwtString))
	}
}
