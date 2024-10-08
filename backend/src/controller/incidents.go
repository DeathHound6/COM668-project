package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateIncident godoc
//
//	@Summary Create an incident
//	@Description Create an incident
//	@Tags Incidents
//	@Security ApiToken
//	@Accept json
//	@Produce json
//	@Success 201
//	@Router /incidents [post]
func CreateIncident() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("location", "")
		ctx.Status(http.StatusCreated)
	}
}
