package controller

import (
	"github.com/gin-gonic/gin"
)

// CreateTeam godoc
//
//	@Summary Create a Team
//	@Description Create a Team
//	@Tags Teams
//	@Security ApiToken
//	@Accept json
//	@Produce json
//	@Failure 400 {object} utility.ErrorResponseSchema
//	@Failure 401 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /teams [post]
func CreateTeam() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// DeleteTeam godoc
//
//	@Summary Delete a Team
//	@Description Delete a Team
//	@Tags Teams
//	@Security ApiToken
//	@Accept json
//	@Produce json
//	@Param team_id path string true "Team ID" format(uuid)
//	@Failure 401 {object} utility.ErrorResponseSchema
//	@Failure 404 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /teams/{team_id} [delete]
func DeleteTeam() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
