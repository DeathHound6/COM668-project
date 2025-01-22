package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"fmt"
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
//	@Success		200	{array}		utility.IncidentGetResponseBodySchema
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/incidents [get]
func GetIncidents() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resolved := ctx.Query("resolved")
		resolvedBool, err := strconv.ParseBool(resolved)
		if err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		incidents, err := database.GetIncidents(ctx, &resolvedBool)
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		response := make([]utility.IncidentGetResponseBodySchema, len(incidents))
		for i, incident := range incidents {
			response[i] = utility.IncidentGetResponseBodySchema{
				UUID: incident.UUID,
			}
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
