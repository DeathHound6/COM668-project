package controller

import (
	"com668-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(engine *gin.Engine) {
	// Register Authentication endpoints
	engine.GET(
		"/authorise/slack",
		middleware.TransactionRequestMW(),
		SlackRedirect(),
		middleware.TransactionResponseMW(),
		middleware.ErrorHandlerResponseMW(),
	)
	engine.GET(
		"/authorise/slack/callback",
		middleware.TransactionRequestMW(),
		AuthoriseSlack(),
		middleware.TransactionResponseMW(),
		middleware.ErrorHandlerResponseMW(),
	)

	engine.GET(
		"/authorise/teams",
		middleware.TransactionRequestMW(),
		// TeamsRedirect(),
		middleware.TransactionResponseMW(),
		middleware.ErrorHandlerResponseMW(),
	)
	engine.GET(
		"/authorise/teams/callback",
		middleware.TransactionRequestMW(),
		// AuthoriseTeams(),
		middleware.TransactionResponseMW(),
		middleware.ErrorHandlerResponseMW(),
	)

	// Register users endpoints
	engine.POST(
		"/users",
		middleware.TransactionRequestMW(),
		CreateUser(),
		middleware.TransactionResponseMW(),
		middleware.ErrorHandlerResponseMW(),
	)
	engine.POST(
		"/users/login",
		middleware.TransactionRequestMW(),
		LoginUser(),
		middleware.TransactionResponseMW(),
		middleware.ErrorHandlerResponseMW(),
	)
}
