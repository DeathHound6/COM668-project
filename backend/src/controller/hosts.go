package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetManyHostsResponseSchema utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]

// GetHosts godoc
//
//	@Summary		Get a list of Hosts
//	@Description	Get a list of Hosts
//	@Tags			Hosts
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false		"Page number"
//	@Param			pageSize	query		int	false		"Number of items per page"
//	@Param			hostnames	query		string	false	"Server hostname"
//	@Success		200			{object}	GetManyHostsResponseSchema
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

		filters := database.GetHostsFilters{
			Page:     &page,
			PageSize: &pageSize,
		}

		if hostname, ok := ctx.GetQuery("hostnames"); ok && len(hostname) > 0 {
			hostnames := make([]string, 0)
			hostnames = append(hostnames, strings.Split(hostname, ",")...)
			filters.Hostnames = &hostnames
		}

		hosts, count, err := database.GetHosts(ctx, filters)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		resp := &utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]{
			Data: make([]*utility.HostMachineGetResponseBodySchema, 0),
			Meta: utility.MetaSchema{
				TotalItems: count,
				Pages:      int(math.Ceil(float64(count) / float64(pageSize))),
				Page:       page,
				PageSize:   pageSize,
			},
		}
		for _, host := range hosts {
			users := make([]utility.UserGetResponseBodySchema, 0)
			for _, user := range host.Team.Users {
				users = append(users, utility.UserGetResponseBodySchema{
					UUID:    user.UUID,
					Name:    user.Name,
					SlackID: user.SlackID,
					Admin:   &user.Admin,
				})
			}
			h := &utility.HostMachineGetResponseBodySchema{
				UUID:     host.UUID,
				Hostname: host.Hostname,
				IP4:      host.IP4,
				IP6:      host.IP6,
				OS:       host.OS,
				Team: utility.TeamGetResponseBodySchema{
					UUID:  host.Team.UUID,
					Name:  host.Team.Name,
					Users: users,
				},
			}
			resp.Data = append(resp.Data, h)
		}
		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", resp)
	}
}

// GetHost godoc
//
//	@Summary		Get a Host
//	@Description	Get a Host
//	@Tags			Hosts
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			host_id	path		string	true	"Host UUID"
//	@Success		200		{object}	utility.HostMachineGetResponseBodySchema
//	@Failure		400		{object}	utility.ErrorResponseSchema
//	@Failure		401		{object}	utility.ErrorResponseSchema
//	@Failure		403		{object}	utility.ErrorResponseSchema
//	@Failure		404		{object}	utility.ErrorResponseSchema
//	@Failure		500		{object}	utility.ErrorResponseSchema
//	@Router			/hosts/{host_id} [get]
func GetHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostUUID := ctx.Param("host_id")
		if _, err := uuid.Parse(hostUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "Invalid Host UUID",
			})
			ctx.Next()
			return
		}

		host, err := database.GetHost(ctx, database.GetHostsFilters{
			UUIDs: []string{hostUUID},
		})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

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
		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", h)
	}
}

// CreateHost godoc
//
//	@Summary		Create a Host
//	@Description	Create a Host
//	@Tags			Hosts
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			body	body	utility.HostMachinePostPutRequestBodySchema	true	"Host creation request"
//	@Header			201		header	string										"GET URL"
//	@Success		201
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/hosts [post]
func CreateHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body *utility.HostMachinePostPutRequestBodySchema
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		if status, err := body.Validate(); err != nil {
			ctx.Set("Status", status)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		team, err := database.GetTeam(ctx, database.GetTeamsFilters{
			UUIDs: []string{body.TeamID},
		})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		host := &database.HostMachine{
			OS:       body.OS,
			Hostname: body.Hostname,
			IP4:      body.IP4,
			IP6:      body.IP6,
			TeamID:   team.ID,
		}
		if err := database.CreateHost(ctx, host); err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		ctx.Header("Location", fmt.Sprintf("%s://%s/hosts/%s", ctx.Request.URL.Scheme, ctx.Request.URL.Host, host.UUID))
		ctx.Set("Status", http.StatusCreated)
	}
}

// UpdateHost godoc
//
//	@Summary		Update a Host
//	@Description	Update a Host
//	@Tags			Hosts
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			host_id	path	string										true	"Host UUID"
//	@Param			body	body	utility.HostMachinePostPutRequestBodySchema	true	"Host update request"
//	@Success		204
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		404	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/hosts/{host_id} [put]
func UpdateHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostUUID := ctx.Param("host_id")
		if _, err := uuid.Parse(hostUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "Invalid Host UUID",
			})
			ctx.Next()
			return
		}

		var body *utility.HostMachinePostPutRequestBodySchema
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		if status, err := body.Validate(); err != nil {
			ctx.Set("Status", status)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		host, err := database.GetHost(ctx, database.GetHostsFilters{
			UUIDs: []string{hostUUID},
		})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		team, err := database.GetTeam(ctx, database.GetTeamsFilters{
			UUIDs: []string{body.TeamID},
		})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		host.TeamID = team.ID
		host.OS = body.OS
		host.Hostname = body.Hostname
		host.IP4 = body.IP4
		host.IP6 = body.IP6

		if err := database.UpdateHost(ctx, host); err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		ctx.Set("Status", http.StatusNoContent)
	}
}

// DeleteHost godoc
//
//	@Summary		Delete a Host
//	@Description	Delete a Host
//	@Tags			Hosts
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			host_id	path	string	true	"Host UUID"
//	@Success		204
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		404	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/hosts/{host_id} [delete]
func DeleteHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostUUID := ctx.Param("host_id")
		if _, err := uuid.Parse(hostUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "Invalid Host UUID",
			})
			ctx.Next()
			return
		}

		_, err := database.GetHost(ctx, database.GetHostsFilters{
			UUIDs: []string{hostUUID},
		})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		if err := database.DeleteHost(ctx, hostUUID); err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		ctx.Set("Status", http.StatusNoContent)
	}
}
