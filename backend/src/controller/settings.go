package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetProviders godoc
//
//	@Summary Get a list of Providers
//	@Description Get a list of Providers
//	@Tags Settings
//	@Security JWT
//	@Accept json
//	@Produce json
//	@Param provider_type query string true "The type of provider" Enums(log, alert)
//	@Success 200 {object} utility.ProvidersGetResponseSchema
//	@Failure 401 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers [get]
func GetProviders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerType := strings.ToLower(ctx.Query("provider_type"))
		if providerType == "log" {
			providers, err := database.GetLogProviders(ctx, map[string]any{})
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			resp := &utility.ProvidersGetResponseSchema{
				Providers: make([]utility.ProviderGetResponseSchema, 0),
			}
			for _, provider := range providers {
				prov := utility.ProviderGetResponseSchema{
					ID:     provider.UUID,
					Name:   provider.Name,
					Fields: utility.GetFieldsMapFromString(provider.Fields),
				}
				resp.Providers = append(resp.Providers, prov)
			}
			ctx.Set("Status", http.StatusOK)
			ctx.Set("Body", resp)
		} else if providerType == "alert" {
			providers, err := database.GetAlertProviders(ctx, map[string]any{})
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			resp := &utility.ProvidersGetResponseSchema{
				Providers: make([]utility.ProviderGetResponseSchema, 0),
			}
			for _, provider := range providers {
				prov := utility.ProviderGetResponseSchema{
					ID:     provider.UUID,
					Name:   provider.Name,
					Fields: utility.GetFieldsMapFromString(provider.Fields),
				}
				resp.Providers = append(resp.Providers, prov)
			}
			ctx.Set("Status", http.StatusOK)
			ctx.Set("Body", resp)
		} else {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
			ctx.Next()
			return
		}
	}
}

// CreateProvider godoc
//
//	@Summary Create a provider
//	@Description Create a provider
//	@Tags Settings
//	@Security JWT
//	@Accept json
//	@Produce json
//	@Success 204
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers [post]
func CreateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// UpdateProvider godoc
//
//	@Summary Update a provider
//	@Description Update a provider
//	@Tags Settings
//	@Security JWT
//	@Accept json
//	@Produce json
//	@Param provider_id path string true "Provider ID" format(uuid)
//	@Success 204
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers/{provider_id} [put]
func UpdateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// GetSettings godoc
//
//	@Summary Get a list of Settings for a given Provider
//	@Description Get a list of Settings for a given Provider
//	@Tags Settings
//	@Security JWT
//	@Accept json
//	@Produce json
//	@Param provider_type query string true "The type of provider" Enums(log, alert)
//	@Param provider_id path string true "Provider ID" format(uuid)
//	@Success 200 {object} utility.SettingsGetResponseSchema
//	@Failure 401 {object} utility.ErrorResponseSchema
//	@Failure 404 {object} utility.ErrorResponseSchema
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers/{provider_id}/settings [get]
func GetSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerType := strings.ToLower(ctx.Query("provider_type"))
		providerUuid := ctx.Param("provider_id")
		if providerType == "log" {
			prov, err := database.GetLogProvider(ctx, providerUuid)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			if prov == nil {
				ctx.Set("Status", http.StatusNotFound)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: fmt.Sprintf("provider with uuid '%s' not found", providerUuid),
				})
				ctx.Next()
				return
			}
			settings, err := database.GetLogSettings(ctx, prov.ID)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			if settings == nil {
				ctx.Set("Status", http.StatusNotFound)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: fmt.Sprintf("no settings found for provider '%s'", providerUuid),
				})
				ctx.Next()
				return
			}
			ctx.Set("Status", http.StatusOK)
			ctx.Set("Body", &utility.SettingsGetResponseSchema{
				Provider: settings.Provider.Name,
				Settings: utility.GetFieldsMapFromString(settings.Settings),
			})
		} else if providerType == "alert" {
			prov, err := database.GetAlertProvider(ctx, providerUuid)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			if prov == nil {
				ctx.Set("Status", http.StatusNotFound)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: fmt.Sprintf("provider with uuid '%s' not found", providerUuid),
				})
				ctx.Next()
				return
			}
			settings, err := database.GetAlertSettings(ctx, prov.ID)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			if settings == nil {
				ctx.Set("Status", http.StatusNotFound)
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: fmt.Sprintf("no settings found for provider '%s'", providerUuid),
				})
				ctx.Next()
				return
			}
			ctx.Set("Status", http.StatusOK)
			ctx.Set("Body", &utility.SettingsGetResponseSchema{
				Provider: settings.Provider.Name,
				Settings: utility.GetFieldsMapFromString(settings.Settings),
			})
		} else {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
		}
	}
}

// UpdateSettings godoc
//
//	@Summary Create or Update the Settings for a given Provider
//	@Description Create or Update the Settings for a given Provider
//	@Tags Settings
//	@Security JWT
//	@Accept json
//	@Produce json
//	@Param provider_type query string true "The type of provider" Enums(log, alert)
//	@Param provider_id path string true "Provider ID" format(uuid)
//	@Success 204
//	@Failure 500 {object} utility.ErrorResponseSchema
//	@Router /providers/{provider_id}/settings [put]
func UpdateSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body *utility.SettingsPatchRequestBody
		if err := ctx.BindJSON(&body); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		providerType := strings.ToLower(ctx.Query("provider_type"))
		providerUuid := ctx.Param("provider_id")
		if providerType == "log" {
			prov, err := database.GetLogProvider(ctx, providerUuid)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			// Verify fields exist in provider
			settings, err := database.GetLogSettings(ctx, prov.ID)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			settings.Settings = utility.GetStringFromFieldsMap(body.Settings)
			// TODO
		} else if providerType == "alert" {
			prov, err := database.GetAlertProvider(ctx, providerUuid)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			// Verify fields exist in provider
			settings, err := database.GetAlertSettings(ctx, prov.ID)
			if err != nil {
				ctx.Set("Status", ctx.GetInt("errorCode"))
				ctx.Set("Body", &utility.ErrorResponseSchema{
					Error: err.Error(),
				})
				ctx.Next()
				return
			}
			settings.Settings = utility.GetStringFromFieldsMap(body.Settings)
			// TODO
		} else {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
		}
	}
}
