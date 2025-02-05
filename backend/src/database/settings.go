package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Provider struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID   string `gorm:"column:uuid;size:36;unique;not null"`
	Name   string `gorm:"column:name;size:30;unique;not null"`
	Fields string `gorm:"column:fields;size:200;not null"`
	Type   string `gorm:"column:type;check:type IN ('log','alert');size:5;not null"`
}

func (p *Provider) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create a provider uuid")
	}
	if len(p.Name) > 30 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("provider name cannot be greater than 30 characters")
	}
	p.UUID = uuid
	return nil
}

type GetProvidersFilters struct {
	UUID         *string
	ProviderType *string
	Name         *string
	Page         *int
	PageSize     *int
}

// Get a single Provider by UUID
func GetProvider(ctx *gin.Context, filters GetProvidersFilters) (*Provider, error) {
	providers, count, err := GetProviders(ctx, GetProvidersFilters{
		UUID:         filters.UUID,
		ProviderType: filters.ProviderType,
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		ctx.Set("errorCode", http.StatusNotFound)
		return nil, errors.New("setting not found")
	}
	return providers[0], nil
}

// Get a list of providers
func GetProviders(ctx *gin.Context, filters GetProvidersFilters) ([]*Provider, int64, error) {
	tx := GetDBTransaction(ctx).Model(&Provider{})

	// apply filters
	if filters.UUID != nil {
		tx = tx.Where("uuid = ?", *filters.UUID)
	}
	if filters.ProviderType != nil {
		tx = tx.Where("type = ?", *filters.ProviderType)
	}
	if filters.Name != nil {
		tx = tx.Where("name = ?", *filters.Name)
	}

	var count int64
	tx.Count(&count)
	if filters.PageSize != nil {
		tx = tx.Limit(*filters.PageSize)
		if filters.Page != nil {
			tx = tx.Offset((*filters.Page - 1) * *filters.PageSize)
		}
	}

	providers := make([]*Provider, 0)
	tx = tx.Find(&providers)
	if tx.Error != nil {
		return nil, -1, handleError(ctx, tx.Error)
	}
	return providers, count, nil
}

// Create a Provider
func CreateProvider(ctx *gin.Context, provider *Provider) error {
	tx := GetDBTransaction(ctx).Model(&Provider{})
	tx = tx.Create(provider)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

// Update a Provider
func UpdateProvider(ctx *gin.Context, uuid string, fields string) error {
	tx := GetDBTransaction(ctx).Model(&Provider{})
	tx = tx.Where("uuid = ?", uuid).Update("fields", fields)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

// Delete a provider
func DeleteProvider(ctx *gin.Context, uuid string) error {
	tx := GetDBTransaction(ctx).Model(&Provider{})
	tx = tx.Where("uuid = ?", uuid).Delete(&Provider{})
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}
