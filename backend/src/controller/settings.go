package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetProviders godoc
//
//	@Summary Get a list of Providers
//	@Description Get a list of Providers
//	@Tags Settings
//	@Security ApiToken
//	@Accept json
//	@Produce json
//	@Param provider_type query string true "The type of provider" Enums(log, alert)
//	@Success 200 {object} utility.ProvidersGetResponseSchema
//	@Failure 401 {object} utility.ErrorResponseSchema
//	@Failure 404 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers [get]
func GetProviders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerType := strings.ToLower(ctx.Query("provider_type"))
		if providerType == "log" {
			providers, err := database.GetLogProviders(ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(ctx.GetInt("errorCode"), &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
			}
			resp := &utility.ProvidersGetResponseSchema{
				Providers: make([]utility.ProviderGetResponseSchema, 0),
			}
			for _, provider := range providers {
				prov := utility.ProviderGetResponseSchema{
					ID:       provider.UUID,
					Name:     provider.Name,
					ImageURL: provider.ImageURL,
					Fields:   utility.GetFieldsMapFromString(provider.Fields),
				}
				resp.Providers = append(resp.Providers, prov)
			}
			ctx.JSON(http.StatusOK, resp)
		} else if providerType == "alert" {
			// TODO
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
			ctx.Next()
		}
	}
}

// GetSettings godoc
//
//	@Summary Get a list of Settings for a given Provider
//	@Description Get a list of Settings for a given Provider
//	@Tags Settings
//	@Security ApiToken
//	@Accept json
//	@Produce json
//	@Param provider_type query string true "The type of provider" Enums(log, alert)
//	@Param provider_id path string true "Provider ID" format(uuid)
//	@Failure 401 {object} utility.ErrorResponseSchema
//	@Failure 404 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers/{provider_id}/settings [get]
func GetSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerType := strings.ToLower(ctx.Query("provider_type"))
		// providerUuid := ctx.Param("provider_id")
		if providerType == "log" {
			// TODO
		} else if providerType == "alert" {
			// TODO
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
			ctx.Next()
		}
	}
}

// PatchSettings godoc
//
//	@Summary Create or Update the list of Settings for a given Provider
//	@Description Create or Update the list of Settings for a given Provider
//	@Tags Settings
//	@Security ApiToken
//	@Accept json
//	@Produce json
//	@Param provider_type query string true "The type of provider" Enums(log, alert)
//	@Param provider_id path string true "Provider ID" format(uuid)
//	@Failure 400 {object} utility.ErrorResponseSchema
//	@Failure 401 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers/{provider_id}/settings [patch]
func PatchSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerType := strings.ToLower(ctx.Query("provider_type"))
		// providerUuid := ctx.Param("provider_id")
		if providerType == "log" {
			// TODO
		} else if providerType == "alert" {
			// TODO
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
			ctx.Next()
		}
	}
}
