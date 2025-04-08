package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Provider struct {
	ID     uint            `gorm:"column:id;primaryKey;autoIncrement"`
	UUID   string          `gorm:"column:uuid;size:36;unique;not null;uniqueIndex"`
	Name   string          `gorm:"column:name;size:30;unique;not null"`
	Fields []ProviderField `gorm:"foreignKey:provider_id;constraint:OnDelete:CASCADE"`
	Type   string          `gorm:"column:type;check:type IN ('log','alert');size:5;not null"`
}
type ProviderField struct {
	ID         uint     `gorm:"column:id;primaryKey;autoIncrement"`
	Provider   Provider `gorm:"foreignKey:provider_id;references:id"`
	ProviderID uint     `gorm:"column:provider_id;not null"`
	Key        string   `gorm:"column:key;size:20;not null"`
	Value      string   `gorm:"column:value;size:30;not null"`
	Type       string   `gorm:"column:type;check:type IN ('string','number','bool');size:10;not null"`
	Required   bool     `gorm:"column:required;not null"`
}

func (p *Provider) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if p.UUID == "" {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusInternalServerError)
			}
			return errors.New("failed to create a provider uuid")
		}
		p.UUID = uuid
	}
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
		PageSize:     utility.Pointer(1),
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
	tx = tx.Preload("Fields")

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
func UpdateProvider(ctx *gin.Context, provider *Provider) error {
	tx := GetDBTransaction(ctx)

	// update name
	if err := tx.Model(&Provider{}).Where("id = ?", provider.ID).Update("name", provider.Name).Error; err != nil {
		return handleError(ctx, tx.Error)
	}
	// replace fields
	if err := tx.Model(&ProviderField{}).Where("provider_id = ?", provider.ID).Delete(&ProviderField{}).Error; err != nil {
		return handleError(ctx, err)
	}
	fields := make([]ProviderField, 0)
	for _, field := range provider.Fields {
		providerField := ProviderField{
			ProviderID: provider.ID,
			Key:        field.Key,
			Value:      field.Value,
			Type:       field.Type,
			Required:   field.Required,
		}
		fields = append(fields, providerField)
	}
	if err := tx.Model(&ProviderField{}).Create(&fields).Error; err != nil {
		return handleError(ctx, err)
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
