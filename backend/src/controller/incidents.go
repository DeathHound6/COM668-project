package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIncidents godoc
//
//	@Summary		Get a list of incidents
//	@Description	Get a list of incidents
//	@Tags			Incidents
//	@Security		JWT
//	@Produce		json
//	@Param			page		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Number of items per page"
//	@Param			resolved	query		bool	false	"Filter by resolved status"
//	@Success		200			{object}	utility.GetManyResponseSchema
//	@Failure		401			{object}	utility.ErrorResponseSchema
//	@Failure		500			{object}	utility.ErrorResponseSchema
//	@Router			/incidents [get]
func GetIncidents() gin.HandlerFunc {
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

		resolved := ctx.Query("resolved")
		filters := database.GetIncidentsFilters{
			Page:     &page,
			PageSize: &pageSize,
		}
		if resolved != "" {
			resolvedBool, err := strconv.ParseBool(resolved)
			if err != nil {
				ctx.Set("Status", http.StatusBadRequest)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			filters.Resolved = &resolvedBool
		}

		incidents, count, err := database.GetIncidents(ctx, filters)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		response := &utility.GetManyResponseSchema[*utility.IncidentGetResponseBodySchema]{
			Data: make([]*utility.IncidentGetResponseBodySchema, 0),
			Meta: utility.MetaSchema{
				TotalItems: count,
				Pages:      int(math.Ceil(float64(count) / float64(pageSize))),
				Page:       page,
				PageSize:   pageSize,
			},
		}
		for _, incident := range incidents {
			inc := &utility.IncidentGetResponseBodySchema{
				UUID:          incident.UUID,
				HostsAffected: make([]utility.HostMachineGetResponseBodySchema, 0),
				Summary:       incident.Summary,
				ResolvedAt:    incident.ResolvedAt,
				CreatedAt:     incident.CreatedAt,
				ResolvedBy: &utility.UserGetResponseBodySchema{
					UUID:    incident.ResolvedBy.UUID,
					Name:    incident.ResolvedBy.Name,
					Email:   incident.ResolvedBy.Email,
					Teams:   make([]utility.TeamGetResponseBodySchema, 0),
					SlackID: incident.ResolvedBy.SlackID,
					Admin:   incident.ResolvedBy.Admin,
				},
			}
			for _, host := range incident.HostsAffected {
				inc.HostsAffected = append(inc.HostsAffected, utility.HostMachineGetResponseBodySchema{
					UUID:     host.UUID,
					Hostname: host.Hostname,
					OS:       host.OS,
					IP4:      host.IP4,
					IP6:      host.IP6,
					Team: utility.TeamGetResponseBodySchema{
						UUID: host.Team.UUID,
						Name: host.Team.Name,
					},
				})
			}
			for _, team := range incident.ResolvedBy.Teams {
				inc.ResolvedBy.Teams = append(inc.ResolvedBy.Teams, utility.TeamGetResponseBodySchema{
					UUID: team.UUID,
					Name: team.Name,
				})
			}
			response.Data = append(response.Data, inc)
		}

		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", response)
	}
}

// CreateIncident godoc
//
//	@Summary		Create an incident
//	@Description	Create an incident
//	@Tags			Incidents
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Header			201	header	string	"GET URL"
//	@Success		201
//	@Router			/incidents [post]
func CreateIncident() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body *utility.IncidentPostRequestBodySchema
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		incident, err := database.CreateIncident(ctx, body)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		ctx.Header("location", fmt.Sprintf("%s://%s/incidents/%s", ctx.Request.URL.Scheme, ctx.Request.URL.Host, incident.UUID))
		ctx.Set("Status", http.StatusCreated)
	}
}
