package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetManyIncidentsResponseSchema utility.GetManyResponseSchema[*utility.IncidentGetResponseBodySchema]

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
//	@Param			myTeams		query		bool	false	"Filter by my teams only"
//	@Param			hash		query		string	false	"Filter by hash"
//	@Success		200			{object}	GetManyIncidentsResponseSchema
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

		filters := database.GetIncidentsFilters{
			Page:     &page,
			PageSize: &pageSize,
		}
		if resolved := ctx.Query("resolved"); resolved != "" {
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

		myTeams := ctx.Query("myTeams")
		if myTeams == "" {
			myTeams = "false"
		}
		myTeamsBool, err := strconv.ParseBool(myTeams)
		if err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		filters.MyTeams = myTeamsBool

		hash := ctx.Query("hash")
		if hash != "" {
			filters.Hash = &hash
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
				UUID:            incident.UUID,
				Comments:        make([]utility.IncidentCommentGetResponseBodySchema, 0),
				HostsAffected:   make([]utility.HostMachineGetResponseBodySchema, 0),
				Description:     incident.Description,
				Summary:         incident.Summary,
				ResolvedAt:      incident.ResolvedAt,
				CreatedAt:       incident.CreatedAt,
				ResolutionTeams: make([]utility.TeamGetResponseBodySchema, 0),
				Hash:            incident.Hash,
			}
			for _, team := range incident.ResolutionTeams {
				users := make([]utility.UserGetResponseBodySchema, 0)
				for _, user := range team.Users {
					users = append(users, utility.UserGetResponseBodySchema{
						UUID:    user.UUID,
						Name:    user.Name,
						Email:   user.Email,
						SlackID: user.SlackID,
						Admin:   &user.Admin,
					})
				}
				inc.ResolutionTeams = append(inc.ResolutionTeams, utility.TeamGetResponseBodySchema{
					UUID:  team.UUID,
					Name:  team.Name,
					Users: users,
				})
			}
			for _, comment := range incident.Comments {
				c := utility.IncidentCommentGetResponseBodySchema{
					UUID:        comment.UUID,
					Comment:     comment.Comment,
					CommentedAt: comment.CommentedAt,
					CommentedBy: utility.UserGetResponseBodySchema{
						UUID:    comment.CommentedBy.UUID,
						Name:    comment.CommentedBy.Name,
						Email:   comment.CommentedBy.Email,
						Teams:   make([]utility.TeamGetResponseBodySchema, 0),
						SlackID: comment.CommentedBy.SlackID,
						Admin:   &comment.CommentedBy.Admin,
					},
				}
				for _, team := range comment.CommentedBy.Teams {
					c.CommentedBy.Teams = append(c.CommentedBy.Teams, utility.TeamGetResponseBodySchema{
						UUID: team.UUID,
						Name: team.Name,
					})
				}
				inc.Comments = append(inc.Comments, c)
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
			// check for nullptr
			if incident.ResolvedBy != nil {
				inc.ResolvedBy = &utility.UserGetResponseBodySchema{
					UUID:    incident.ResolvedBy.UUID,
					Name:    incident.ResolvedBy.Name,
					Email:   incident.ResolvedBy.Email,
					Teams:   make([]utility.TeamGetResponseBodySchema, 0),
					SlackID: incident.ResolvedBy.SlackID,
					Admin:   &incident.ResolvedBy.Admin,
				}
				for _, team := range incident.ResolvedBy.Teams {
					inc.ResolvedBy.Teams = append(inc.ResolvedBy.Teams, utility.TeamGetResponseBodySchema{
						UUID: team.UUID,
						Name: team.Name,
					})
				}
			}
			response.Data = append(response.Data, inc)
		}

		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", response)
	}
}

// GetIncident godoc
//
//	@Summary		Get an incident
//	@Description	Get an incident
//	@Tags			Incidents
//	@Security		JWT
//	@Produce		json
//	@Param			incident_id	path		string	true	"Incident UUID"
//	@Success		200			{object}	utility.IncidentGetResponseBodySchema
//	@Failure		401			{object}	utility.ErrorResponseSchema
//	@Failure		404			{object}	utility.ErrorResponseSchema
//	@Failure		500			{object}	utility.ErrorResponseSchema
//	@Router			/incidents/{incident_id} [get]
func GetIncident() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		incidentUUID := ctx.Param("incident_id")
		if _, err := uuid.Parse(incidentUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "invalid incident UUID",
			})
			ctx.Next()
			return
		}

		incident, err := database.GetIncident(ctx, incidentUUID)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		var resolvedBy *utility.UserGetResponseBodySchema = nil
		if incident.ResolvedBy != nil {
			resolvedBy = &utility.UserGetResponseBodySchema{
				UUID:    incident.ResolvedBy.UUID,
				Name:    incident.ResolvedBy.Name,
				Email:   incident.ResolvedBy.Email,
				Teams:   make([]utility.TeamGetResponseBodySchema, 0),
				SlackID: incident.ResolvedBy.SlackID,
				Admin:   &incident.ResolvedBy.Admin,
			}
			for _, team := range incident.ResolvedBy.Teams {
				resolvedBy.Teams = append(resolvedBy.Teams, utility.TeamGetResponseBodySchema{
					UUID: team.UUID,
					Name: team.Name,
				})
			}
		}
		inc := &utility.IncidentGetResponseBodySchema{
			UUID:            incident.UUID,
			Comments:        make([]utility.IncidentCommentGetResponseBodySchema, 0),
			HostsAffected:   make([]utility.HostMachineGetResponseBodySchema, 0),
			Description:     incident.Description,
			Summary:         incident.Summary,
			ResolvedAt:      incident.ResolvedAt,
			ResolvedBy:      resolvedBy,
			CreatedAt:       incident.CreatedAt,
			ResolutionTeams: make([]utility.TeamGetResponseBodySchema, 0),
			Hash:            incident.Hash,
		}
		for _, team := range incident.ResolutionTeams {
			users := make([]utility.UserGetResponseBodySchema, 0)
			for _, user := range team.Users {
				users = append(users, utility.UserGetResponseBodySchema{
					UUID:    user.UUID,
					Name:    user.Name,
					SlackID: user.SlackID,
					Admin:   &user.Admin,
				})
			}
			inc.ResolutionTeams = append(inc.ResolutionTeams, utility.TeamGetResponseBodySchema{
				UUID:  team.UUID,
				Name:  team.Name,
				Users: users,
			})
		}
		for _, comment := range incident.Comments {
			c := utility.IncidentCommentGetResponseBodySchema{
				UUID:        comment.UUID,
				Comment:     comment.Comment,
				CommentedAt: comment.CommentedAt,
				CommentedBy: utility.UserGetResponseBodySchema{
					UUID:    comment.CommentedBy.UUID,
					Name:    comment.CommentedBy.Name,
					Email:   comment.CommentedBy.Email,
					Teams:   make([]utility.TeamGetResponseBodySchema, 0),
					SlackID: comment.CommentedBy.SlackID,
					Admin:   &comment.CommentedBy.Admin,
				},
			}
			for _, team := range comment.CommentedBy.Teams {
				c.CommentedBy.Teams = append(c.CommentedBy.Teams, utility.TeamGetResponseBodySchema{
					UUID: team.UUID,
					Name: team.Name,
				})
			}
			inc.Comments = append(inc.Comments, c)
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

		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", inc)
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
//	@Param			incident	body	utility.IncidentPostRequestBodySchema	true	"The request body"
//	@Header			201			header	string									"GET URL"
//	@Success		201
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
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

		if status, err := body.Validate(); err != nil {
			ctx.Set("Status", status)
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

		ctx.Header("Location", fmt.Sprintf("%s://%s/incidents/%s", ctx.Request.URL.Scheme, ctx.Request.URL.Host, incident.UUID))
		ctx.Set("Status", http.StatusCreated)
	}
}

// UpdateIncident godoc
//
//	@Summary		Update an incident
//	@Description	Update an incident
//	@Tags			Incidents
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			incident	body	utility.IncidentPutRequestBodySchema	true	"The request body"
//	@Param			incident_id	path	string									true	"Incident UUID"
//	@Success		204
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		404	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/incidents/{incident_id} [put]
func UpdateIncident() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		incidentUUID := ctx.Param("incident_id")
		if _, err := uuid.Parse(incidentUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "invalid incident UUID",
			})
			ctx.Next()
			return
		}

		var body *utility.IncidentPutRequestBodySchema
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

		incident, err := database.GetIncident(ctx, incidentUUID)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		hosts := make([]database.HostMachine, 0)
		if len(body.HostsAffected) > 0 {
			hs, count, err := database.GetHosts(ctx, database.GetHostsFilters{
				UUIDs:    body.HostsAffected,
				PageSize: utility.Pointer(len(body.HostsAffected)),
			})
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			if int(count) != len(body.HostsAffected) {
				ctx.Set("Status", http.StatusBadRequest)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: "one or more hosts not found",
				})
				ctx.Next()
				return
			}
			for _, host := range hs {
				hosts = append(hosts, *host)
			}
		}

		teams := make([]database.Team, 0)
		if len(body.ResolutionTeams) > 0 {
			ts, count, err := database.GetTeams(ctx, database.GetTeamsFilters{
				UUIDs:    body.ResolutionTeams,
				PageSize: utility.Pointer(len(body.ResolutionTeams)),
			})
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			if int(count) != len(body.ResolutionTeams) {
				ctx.Set("Status", http.StatusBadRequest)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: "one or more teams not found",
				})
				ctx.Next()
				return
			}
			for _, team := range ts {
				teams = append(teams, *team)
			}
		}

		var resolvedByID *uint = nil
		var resolvedAt *time.Time = nil
		if body.Resolved != nil && *body.Resolved {
			user := ctx.MustGet("user").(*database.User)
			resolvedByID = &user.ID
			resolvedAt = utility.Pointer(time.Now())
		}

		newIncident := &database.Incident{
			ID:              incident.ID,
			UUID:            incident.UUID,
			Summary:         body.Summary,
			Description:     body.Description,
			HostsAffected:   hosts,
			Comments:        incident.Comments,
			ResolvedByID:    resolvedByID,
			ResolvedAt:      resolvedAt,
			CreatedAt:       incident.CreatedAt,
			ResolutionTeams: teams,
		}
		err = database.UpdateIncident(ctx, database.GetIncidentsFilters{
			UUID: &incidentUUID,
		}, newIncident)
		if err != nil {
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

// CreateIncidentComment godoc
//
//	@Summary		Create an incident comment
//	@Description	Create an incident comment
//	@Tags			Incidents
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			comment		body	utility.IncidentCommentPostRequestBodySchema	true	"The request body"
//	@Param			incident_id	path	string											true	"Incident UUID"
//	@Success		201
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/incidents/{incident_id}/comments [post]
func CreateIncidentComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		incidentUUID := ctx.Param("incident_id")
		if _, err := uuid.Parse(incidentUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "invalid incident UUID",
			})
			ctx.Next()
			return
		}

		var body *utility.IncidentCommentPostRequestBodySchema
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

		incident, err := database.GetIncident(ctx, incidentUUID)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		user := ctx.MustGet("user").(*database.User)
		comment := &database.IncidentComment{
			Comment:       body.Comment,
			IncidentID:    incident.ID,
			CommentedByID: user.ID,
		}
		comment, err = database.CreateIncidentComment(ctx, comment)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		ctx.Set("Status", http.StatusNoContent)
		ctx.Header("Location", fmt.Sprintf("%s://%s/incidents/%s/comments/%s", ctx.Request.URL.Scheme, ctx.Request.URL.Host, incident.UUID, comment.UUID))
	}
}

// DeleteIncidentComment godoc
//
//	@Summary		Delete an incident comment
//	@Description	Delete an incident comment
//	@Tags			Incidents
//	@Security		JWT
//	@Produce		json
//	@Param			comment_id	path	string	true	"Comment UUID"
//	@Param			incident_id	path	string	true	"Incident UUID"
//	@Success		204
//	@Failure		400	{object}	utility.ErrorResponseSchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		404	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/incidents/{incident_id}/comments/{comment_id} [delete]
func DeleteIncidentComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*database.User)
		incidentUUID := ctx.Param("incident_id")
		if _, err := uuid.Parse(incidentUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "invalid incident UUID",
			})
			ctx.Next()
			return
		}
		commentUUID := ctx.Param("comment_id")
		if _, err := uuid.Parse(commentUUID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "invalid comment UUID",
			})
			ctx.Next()
			return
		}

		comment, err := database.GetIncidentComment(ctx, incidentUUID, commentUUID)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		if !user.Admin && comment.CommentedByID != user.ID {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "you are not allowed to delete this comment",
			})
			ctx.Next()
			return
		}

		err = database.DeleteIncidentComment(ctx, comment.UUID)
		if err != nil {
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
