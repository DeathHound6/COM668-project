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
//	@Summary		Get a list of Providers
//	@Description	Get a list of Providers
//	@Tags			Settings
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			provider_type	query		string	true	"The type of provider"	Enums(log, alert)
//	@Success		200				{object}	utility.ProvidersGetResponseSchema
//	@Failure		401				{object}	utility.ErrorResponseSchema
//	@Failure		403				{object}	utility.ErrorResponseSchema
//	@Failure		500				{object}	utility.ErrorResponseSchema
//	@Router			/providers [get]
func GetProviders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loggedInUser := ctx.MustGet("user").(*database.User)
		if !loggedInUser.Admin {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "logged in user must be an admin",
			})
			ctx.Next()
			return
		}

		providerType := strings.ToLower(ctx.Query("provider_type"))
		if !utility.SliceHasElement([]string{"alert", "log"}, providerType) {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
			ctx.Next()
			return
		}

		filters := make(map[string]any)
		filters["type"] = providerType
		providers, err := database.GetProviders(ctx, filters)
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
				Type:   provider.Type,
			}
			resp.Providers = append(resp.Providers, prov)
		}
		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", resp)
	}
}

// CreateProvider godoc
//
//	@Summary		Create a provider
//	@Description	Create a provider
//	@Tags			Settings
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/providers [post]
func CreateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loggedInUser := ctx.MustGet("user").(*database.User)
		if !loggedInUser.Admin {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "logged in user must be an admin",
			})
			ctx.Next()
			return
		}
	}
}

// UpdateProvider godoc
//
//	@Summary		Update a provider
//	@Description	Update a provider
//	@Tags			Settings
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			provider_id	path	string	true	"Provider ID"	format(uuid)
//	@Success		204
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/providers/{provider_id} [put]
func UpdateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loggedInUser := ctx.MustGet("user").(*database.User)
		if !loggedInUser.Admin {
			ctx.Set("Status", http.StatusForbidden)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "logged in user must be an admin",
			})
			ctx.Next()
			return
		}
	}
}
