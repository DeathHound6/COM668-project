package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetManyTeamsResponseSchema utility.GetManyResponseSchema[*utility.TeamGetResponseBodySchema]

// GetTeams godoc
//
//	@Summary		Get a list of Teams
//	@Description	Get a list of Teams
//	@Tags			Teams
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Number of items per page"
//	@Success		200			{object}	GetManyTeamsResponseSchema
//	@Failure		400			{object}	utility.ErrorResponseSchema
//	@Failure		401			{object}	utility.ErrorResponseSchema
//	@Failure		403			{object}	utility.ErrorResponseSchema
//	@Failure		500			{object}	utility.ErrorResponseSchema
//	@Router			/teams [get]
func GetTeams() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params, err := getCommonParams(ctx)
		if err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		page := params["page"].(int)
		pageSize := params["pageSize"].(int)

		teams, count, err := database.GetTeams(ctx, database.GetTeamsFilters{
			Page:     &page,
			PageSize: &pageSize,
		})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		resp := &utility.GetManyResponseSchema[*utility.TeamGetResponseBodySchema]{
			Data: make([]*utility.TeamGetResponseBodySchema, 0),
			Meta: utility.MetaSchema{
				TotalItems: count,
				Pages:      int(math.Ceil(float64(count) / float64(pageSize))),
				Page:       page,
				PageSize:   pageSize,
			},
		}
		for _, team := range teams {
			t := &utility.TeamGetResponseBodySchema{
				UUID: team.UUID,
				Name: team.Name,
			}
			resp.Data = append(resp.Data, t)
		}
		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", resp)
	}
}

// CreateTeam godoc
//
//	@Summary		Create a Team
//	@Description	Create a Team
//	@Tags			Teams
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/teams [post]
func CreateTeam() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// DeleteTeam godoc
//
//	@Summary		Delete a Team
//	@Description	Delete a Team
//	@Tags			Teams
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			team_id	path		string	true	"Team ID"	format(uuid)
//	@Failure		401		{object}	utility.ErrorResponseSchema
//	@Failure		404		{object}	utility.ErrorResponseSchema
//	@Failure		500		{object}	utility.ErrorResponseSchema
//	@Router			/teams/{team_id} [delete]
func DeleteTeam() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
