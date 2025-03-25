package controller

import (
	"com668-backend/database"
	"com668-backend/utility"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetManyProvidersResponseSchema utility.GetManyResponseSchema[*utility.ProviderGetResponseSchema]

// GetProviders godoc
//
//	@Summary		Get a list of Providers
//	@Description	Get a list of Providers
//	@Tags			Settings
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			provider_type	query		string	true	"The type of provider"	Enums(log, alert)
//	@Param			page			query		int		false	"Page number"
//	@Param			pageSize		query		int		false	"Number of items per page"
//	@Success		200				{object}	GetManyProvidersResponseSchema
//	@Failure		401				{object}	utility.ErrorResponseSchema
//	@Failure		403				{object}	utility.ErrorResponseSchema
//	@Failure		500				{object}	utility.ErrorResponseSchema
//	@Router			/providers [get]
func GetProviders() gin.HandlerFunc {
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

		providerType := strings.ToLower(ctx.Query("provider_type"))
		if !slices.Contains([]string{"alert", "log"}, providerType) {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
			ctx.Next()
			return
		}

		providers, count, err := database.GetProviders(ctx, database.GetProvidersFilters{
			ProviderType: &providerType,
		})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		resp := &utility.GetManyResponseSchema[*utility.ProviderGetResponseSchema]{
			Data: make([]*utility.ProviderGetResponseSchema, 0),
			Meta: utility.MetaSchema{
				Page:       page,
				PageSize:   pageSize,
				TotalItems: count,
				Pages:      int(math.Ceil(float64(count) / float64(pageSize))),
			},
		}
		for _, provider := range providers {
			fmt.Printf("fields %v\n", provider.Fields)
			fields := make([]utility.KeyValueSchema, 0)
			for _, field := range provider.Fields {
				fields = append(fields, utility.KeyValueSchema{
					Key:      field.Key,
					Value:    field.Value,
					Type:     field.Type,
					Required: &field.Required,
				})
			}
			prov := &utility.ProviderGetResponseSchema{
				UUID:   provider.UUID,
				Name:   provider.Name,
				Fields: fields,
				Type:   provider.Type,
			}
			resp.Data = append(resp.Data, prov)
		}
		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", resp)
	}
}

// GetProvider godoc
//
//	@Summary		Get a provider
//	@Description	Get a provider
//	@Tags			Settings
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			provider_id	path		string	true	"Provider ID"	format(uuid)
//	@Success		200			{object}	utility.ProviderGetResponseSchema
//	@Failure		401			{object}	utility.ErrorResponseSchema
//	@Failure		404			{object}	utility.ErrorResponseSchema
//	@Failure		500			{object}	utility.ErrorResponseSchema
//	@Router			/providers/{provider_id} [get]
func GetProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerID := ctx.Param("provider_id")
		if _, err := uuid.Parse(providerID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "Invalid provider ID",
			})
			ctx.Next()
			return
		}

		provider, err := database.GetProvider(ctx, database.GetProvidersFilters{UUID: &providerID})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		fields := make([]utility.KeyValueSchema, 0)
		for _, field := range provider.Fields {
			fields = append(fields, utility.KeyValueSchema{
				Key:      field.Key,
				Value:    field.Value,
				Type:     field.Type,
				Required: &field.Required,
			})
		}

		ctx.Set("Status", http.StatusOK)
		ctx.Set("Body", &utility.ProviderGetResponseSchema{
			UUID:   provider.UUID,
			Name:   provider.Name,
			Fields: fields,
			Type:   provider.Type,
		})
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
//	@Param			provider_type	query	string									true	"The type of provider"	Enums(log, alert)
//	@Param			body			body	utility.ProviderPostRequestBodySchema	true	"Provider data"
//	@Success		201
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/providers [post]
func CreateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerType := strings.ToLower(ctx.Query("provider_type"))
		if !slices.Contains([]string{"alert", "log"}, providerType) {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "'provider_type' query parameter must be either 'log' or 'alert'",
			})
			ctx.Next()
			return
		}
		var body *utility.ProviderPostRequestBodySchema
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

		provider := &database.Provider{
			Name:   body.Name,
			Fields: []database.ProviderField{},
			Type:   providerType,
		}
		if err := database.CreateProvider(ctx, provider); err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		ctx.Set("Status", http.StatusCreated)
		ctx.Header("Location", fmt.Sprintf("%s://%s/providers/%s", ctx.Request.URL.Scheme, ctx.Request.URL.Host, provider.UUID))
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
//	@Param			provider_id	path	string									true	"Provider ID"	format(uuid)
//	@Param			body		body	utility.ProviderPutRequestBodySchema	true	"Provider data"
//	@Success		204
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		404	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/providers/{provider_id} [put]
func UpdateProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerID := ctx.Param("provider_id")
		if _, err := uuid.Parse(providerID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "Invalid provider ID",
			})
			ctx.Next()
			return
		}

		var body *utility.ProviderPutRequestBodySchema
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

		provider, err := database.GetProvider(ctx, database.GetProvidersFilters{UUID: &providerID})
		if err != nil {
			ctx.Set("Status", ctx.GetInt("errorCode"))
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		provider.Name = body.Name
		provider.Fields = []database.ProviderField{}
		for _, field := range body.Fields {
			var required bool = false
			if field.Required != nil {
				required = *field.Required
			}
			providerField := database.ProviderField{
				Key:      field.Key,
				Value:    field.Value,
				Type:     field.Type,
				Required: required,
			}
			provider.Fields = append(provider.Fields, providerField)
		}

		if err := database.UpdateProvider(ctx, provider); err != nil {
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

// DeleteProvider godoc
//
//	@Summary		Delete a provider
//	@Description	Delete a provider
//	@Tags			Settings
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			provider_id	path	string	true	"Provider ID"	format(uuid)
//	@Success		204
//	@Failure		401	{object}	utility.ErrorResponseSchema
//	@Failure		403	{object}	utility.ErrorResponseSchema
//	@Failure		404	{object}	utility.ErrorResponseSchema
//	@Failure		500	{object}	utility.ErrorResponseSchema
//	@Router			/providers/{provider_id} [delete]
func DeleteProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerID := ctx.Param("provider_id")
		if _, err := uuid.Parse(providerID); err != nil {
			ctx.Set("Status", http.StatusBadRequest)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: "Invalid provider ID",
			})
			ctx.Next()
			return
		}

		_, err := database.GetProvider(ctx, database.GetProvidersFilters{UUID: &providerID})
		if err != nil {
			ctx.Set("Status", http.StatusNotFound)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}

		if err := database.DeleteProvider(ctx, providerID); err != nil {
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
