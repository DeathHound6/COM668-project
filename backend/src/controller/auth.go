package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/demisto/slack"
	"golang.org/x/oauth2"
)

type state struct {
	auth string
	ts   time.Time
}

var (
	authStateMapper = make(map[string]state, 0)
)

// SlackRedirect godoc
//
//	@Summary		Redirect to Slack auth login
//	@Description	Redirect to Slack auth login
//	@Tags			Third-Party Auth
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/authorise/slack [get]
func SlackRedirect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*database.User)
		b := make([]byte, 10)
		_, err := rand.Read(b)
		if err != nil {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		globalState := state{auth: hex.EncodeToString(b), ts: time.Now()}
		authStateMapper[user.UUID] = globalState
		conf := &oauth2.Config{
			ClientID:     os.Getenv("SLACK_CLIENT_ID"),
			ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
			Scopes:       []string{"users:read.email", "users:read", "users.profile:read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://slack.com/oauth/authorize",
				TokenURL: "https://slack.com/api/oauth.access", // not actually used here
			},
			RedirectURL: "https://localhost:5000/authorise/slack/callback",
		}
		url := conf.AuthCodeURL(globalState.auth)
		ctx.Redirect(302, url)
	}
}

// AuthoriseSlack godoc
//
//	@Summary		Link Slack to user
//	@Description	Link Slack to user
//	@Tags			Third-Party Auth
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/authorise/slack/callback [get]
func AuthoriseSlack() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*database.User)
		slackState := ctx.Request.URL.Query().Get("state")
		code := ctx.Request.URL.Query().Get("code")
		errStr := ctx.Request.URL.Query().Get("error")

		if errStr != "" {
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: errStr,
			})
			ctx.Next()
			return
		}
		if slackState == "" || code == "" {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "missing state or code",
			})
			ctx.Next()
			return
		}
		globalState, ok := authStateMapper[user.UUID]
		if !ok {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "original auth state not present",
			})
			ctx.Next()
			return
		}
		if slackState != globalState.auth {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "state does not match",
			})
			ctx.Next()
			return
		}
		// As an example, we allow only 5 min between requests
		if time.Since(globalState.ts) > time.Minute*5 {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "state is too old",
			})
			ctx.Next()
			return
		}

		token, err := slack.OAuthAccess(os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_CLIENT_SECRET"), code, "")
		if err != nil {
			ctx.Set("Status", http.StatusUnauthorized)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		s, err := slack.New(slack.SetToken(token.AccessToken))
		if err != nil {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		// Get our own user id
		test, err := s.AuthTest()
		if err != nil {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		user.SlackID = test.UserID
		if err := database.UpdateUser(ctx, user); err != nil {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		ctx.Redirect(302, "http://localhost:3000/dashboard")
	}
}
