package controller

import (
	"com668-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(engine *gin.Engine) {
	// Register Authentication endpoints
	auth := engine.Group("/authorise", middleware.TransactionRequestMW())
	{
		auth.GET("/slack", SlackRedirect(), middleware.TransactionResponseMW())
		auth.GET("/slack/callback", AuthoriseSlack(), middleware.TransactionResponseMW())

		// auth.GET("/teams", middleware.TransactionResponseMW())
		// auth.GET("/teams/callback", middleware.TransactionResponseMW())
	}

	// Register users endpoints
	users := engine.Group("/users", middleware.RequestMW...)
	{
		users.POST("/", CreateUser(), middleware.TransactionResponseMW())
		users.POST("/login", LoginUser(), middleware.TransactionResponseMW())
	}
}
