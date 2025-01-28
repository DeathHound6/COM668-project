package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHosts godoc
//
//	@Summary		Get a list of Hosts
//	@Description	Get a list of Hosts
//	@Tags			Hosts
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Number of items per page"
//	@Success		200			{object}	utility.GetManyResponseSchema
//	@Failure		401			{object}	utility.ErrorResponseSchema
//	@Failure		403			{object}	utility.ErrorResponseSchema
//	@Failure		500			{object}	utility.ErrorResponseSchema
//	@Router			/hosts [get]
func GetHosts() gin.HandlerFunc {
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

		hosts, count, err := database.GetHosts(ctx, database.GetHostsFilters{
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

		resp := &utility.GetManyResponseSchema{
			Data: make([]any, 0),
			Meta: utility.MetaSchema{
				TotalItems: count,
				Pages:      int(math.Ceil(float64(count) / float64(pageSize))),
				Page:       page,
				PageSize:   pageSize,
			},
		}
		for _, host := range hosts {
			h := &utility.HostMachineGetResponseBodySchema{
				UUID:     host.UUID,
				Hostname: host.Hostname,
				IP4:      host.IP4,
				IP6:      host.IP6,
				OS:       host.OS,
				Team: utility.TeamGetResponseBodySchema{
					UUID: host.Team.UUID,
					Name: host.Team.Name,
				},
			}
			resp.Data = append(resp.Data, h)
		}
		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", resp)
	}
}
