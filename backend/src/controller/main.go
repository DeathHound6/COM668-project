package controller

import (
	"com668-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(engine *gin.Engine) {
	// Register Authentication endpoints
	engine.GET(
		"/authorise/slack",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		SlackRedirect(),
		middleware.TransactionResponseMW(),
	)
	engine.GET(
		"/authorise/slack/callback",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		AuthoriseSlack(),
		middleware.TransactionResponseMW(),
	)

	// Register teams endpoints
	engine.POST(
		"/teams",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		CreateTeam(),
		middleware.TransactionResponseMW(),
	)
	engine.DELETE(
		"/teams/:team_id",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		DeleteTeam(),
		middleware.TransactionResponseMW(),
	)

	// Register users endpoints
	engine.POST(
		"/users",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		CreateUser(),
		middleware.TransactionResponseMW(),
	)
	engine.POST(
		"/users/login",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		LoginUser(),
		middleware.TransactionResponseMW(),
	)

	// Register incident endpoints
	engine.POST(
		"/incidents",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		CreateIncident(),
		middleware.TransactionResponseMW(),
	)

	// Register settings endpoints
	engine.GET(
		"/providers",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		GetProviders(),
		middleware.TransactionResponseMW(),
	)
	engine.GET(
		"/providers/:provider_id/settings",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		GetSettings(),
		middleware.TransactionResponseMW(),
	)
	engine.PATCH(
		"/providers/:provider_id/settings",
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		PatchSettings(),
		middleware.TransactionResponseMW(),
	)
}
