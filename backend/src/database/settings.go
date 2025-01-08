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
	filter(filters, allowedFilters, tx)
	tx.Find(&providers)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return providers, nil
}
