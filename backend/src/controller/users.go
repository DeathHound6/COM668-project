package controller

import (
	"com668-backend/database"
	"com668-backend/middleware"
	"com668-backend/utility"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GetUser godoc
//
//	@Summary		Get basic details about the currently logged in user
//	@Description	Get basic details about the currently logged in user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Success		200	{object}	utility.UserGetResponseBodySchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Router			/me [get]
func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*database.User)

		teams := make([]utility.TeamGetResponseBodySchema, len(user.Teams))
		for i, team := range user.Teams {
			teams[i] = utility.TeamGetResponseBodySchema{
				UUID: team.UUID,
				Name: team.Name,
			}
		}
		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", &utility.UserGetResponseBodySchema{
			UUID:    user.UUID,
			Name:    user.Name,
			Email:   user.Email,
			Teams:   teams,
			SlackID: user.SlackID,
			Admin:   &user.Admin,
		})
	}
}

// CreateUser godoc
//
//	@Summary		Create a user
//	@Description	Create a user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	utility.UserPostRequestBodySchema	true	"The request body"
//	@Success		201
//	@Header			201	{string}	Location	"GET URL of the created User"
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/users [post]
func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body *utility.UserPostRequestBodySchema
		if err := ctx.ShouldBindJSON(&body); err != nil {
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
//	@Summary		Login as a user
//	@Description	Login as a user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request_body	body	utility.UserLoginRequestBodySchema	true	"Request Body"
//	@Header			204				header	string								"JWT Token"
//	@Success		204
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/users/login [post]
func LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader(middleware.AuthHeaderNameString)
		if authToken != "" {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "user is already authenticated",
			})
			ctx.Next()
			return
		}

		var body *utility.UserLoginRequestBodySchema
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		user, err := database.GetUser(ctx, body.Email)
		if user == nil || err != nil || !user.ValidatePassword(body.Password) {
			if err != nil {
				log.Default().Println(err)
			}
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "invalid email or password",
			})
			ctx.Next()
			return
		}

		token := jwt.New(middleware.JWTSigningMethod)
		claims := jwt.MapClaims{}
		claims["iss"] = "COM668"
		claims["iat"] = jwt.NewNumericDate(time.Now())
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
		claims["sub"] = base64.StdEncoding.EncodeToString([]byte(user.Email))
		token.Claims = claims
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
		ctx.SetCookie(middleware.AuthHeaderNameString, jwtString, int(time.Hour)*24, "/", ctx.Request.URL.Host, true, false)
	}
}
