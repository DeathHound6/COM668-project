package controller

import (
	"github.com/gin-gonic/gin"
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

	}
}
