package middleware

import (
	"com668-backend/database"
	"com668-backend/utility"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	JWTSigningMethod     *jwt.SigningMethodHMAC = jwt.SigningMethodHS256
	AuthHeaderNameString string                 = "Authorization"
)

func UserAuthRequestMW(adminAuth bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authType := "header"
		jwtString := ctx.GetHeader(AuthHeaderNameString)
		// if jwt not in header, check cookies
		if jwtString == "" {
			authType = "cookie"
			jwt, err := ctx.Cookie(AuthHeaderNameString)
			if err != nil {
				jwtString = ""
			}
			jwtString = jwt
		}
		// if jwt not in header or cookies, fail
		if jwtString == "" {
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "no jwt auth token specified",
			})
			ctx.Next()
			return
		}

		// only attempt to parse `Bearer <token>` format if jwt is in header
		if authType == "header" {
			parts := strings.Split(jwtString, " ")
			if len(parts) != 2 || (len(parts) > 0 && strings.ToLower(parts[0]) != "bearer") {
				ctx.Set("Status", http.StatusUnauthorized)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: "jwt auth token does not follow the format `Bearer <token>`",
				})
				ctx.Next()
				return
			}
			jwtString = parts[1]
		}

		token, err := jwt.Parse(
			jwtString,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
			},
			jwt.WithValidMethods([]string{
				JWTSigningMethod.Name,
			}),
		)
		if err != nil {
			log.Default().Println(err)
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "could not parse jwt auth token",
			})
			ctx.Next()
			return
		}
		if !token.Valid {
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "jwt auth token is not valid",
			})
			ctx.Next()
			return
		}

		exp, err := token.Claims.GetExpirationTime()
		if err != nil {
			log.Default().Println(err)
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "could not parse jwt auth claims",
			})
			ctx.Next()
			return
		}
		if time.Now().After(exp.Time) {
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "jwt auth token is no longer valid",
			})
			ctx.Next()
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			log.Default().Println(err)
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "could not parse jwt auth claims",
			})
			ctx.Next()
			return
		}
		subBytes, err := base64.StdEncoding.DecodeString(sub)
		if err != nil {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		sub = string(subBytes)

		user, err := database.GetUser(ctx, database.GetUserFilters{
			UUID: &sub,
		})
		if err != nil {
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		if user == nil {
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "invalid user credentials",
			})
			ctx.Next()
			return
		}
		if adminAuth && !user.Admin {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "logged in user must be an admin",
			})
			ctx.Next()
			return
		}
		ctx.Set("user", user)
	}
}
