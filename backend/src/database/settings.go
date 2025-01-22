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

// Get a single Provider by UUID
func GetProvider(ctx *gin.Context, uuid string, provider_type string) (*Provider, error) {
	filters := make(map[string]any, 0)
	filters["uuid"] = uuid
	filters["type"] = provider_type
	providers, err := GetProviders(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(providers) == 0 {
		return nil, nil
	}
	return providers[0], nil
}

// Get a list of providers
func GetProviders(ctx *gin.Context, filters map[string]any) ([]*Provider, error) {
	// mapped `filterMapField, dbField`
	allowedFilters := [][]string{
		{"id", "id"},
		{"uuid", "uuid"},
		{"name", "name"},
		{"type", "type"},
	}
	tx := GetDBTransaction(ctx)
	providers := make([]*Provider, 0)
	tx = filter(filters, allowedFilters, tx)
	tx = tx.Find(&providers)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return providers, nil
}

// Create a Provider
func CreateProvider(ctx *gin.Context, provider *Provider) error {
	tx := GetDBTransaction(ctx)
	tx = tx.Create(provider)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

// Update a Provider
func UpdateProvider(ctx *gin.Context, uuid string, fields string) error {
	tx := GetDBTransaction(ctx)
	tx = tx.Model(&Provider{}).Where("uuid = ?", uuid).Update("fields", fields)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

// Delete a provider
func DeleteProvider(ctx *gin.Context, uuid string) error {
	tx := GetDBTransaction(ctx)
	tx = tx.Where("uuid = ?", uuid).Delete(&Provider{})
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}
