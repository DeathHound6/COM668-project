package controller

import (
	"github.com/gin-gonic/gin"
)

// @Summary Redirect to Slack auth
// @Description Redirect to Slack auth
// @Tags Third-Party Auth
// @Accept json
// @Produce json
// @Router /authorise/slack [get]
func SlackRedirect() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// @Summary Link Slack to user
// @Description Link Slack to user
// @Tags Third-Party Auth
// @Accept json
// @Produce json
// @Router /authorise/slack/callback [get]
func AuthoriseSlack() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
