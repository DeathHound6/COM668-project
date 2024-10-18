package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogProvider struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID   string `gorm:"column:uuid;size:36;unique;not null"`
	Name   string `gorm:"column:name;size:30;unique;not null"`
	Fields string `gorm:"column:fields;size:200;not null"`
}

type AlertProvider struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID   string `gorm:"column:uuid;size:36;unique;not null"`
	Name   string `gorm:"column:name;size:30;unique;not null"`
	Fields string `gorm:"column:fields;size:200;not null"`
}

func (p *LogProvider) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create a log provider uuid")
	}
	if len(p.Name) > 30 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("log provider name cannot be greater than 30 characters")
	}
	p.UUID = uuid
	return nil
}
func (p *AlertProvider) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create an alert provider uuid")
	}
	if len(p.Name) > 30 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("alert provider name cannot be greater than 30 characters")
	}
	p.UUID = uuid
	return nil
}

// Get a single Log Provider by UUID
func GetLogProvider(ctx *gin.Context, uuid string) (*LogProvider, error) {
	filters := make(map[string]any, 0)
	filters["uuid"] = uuid
	providers, err := GetLogProviders(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(providers) == 0 {
		return nil, nil
	}
	return providers[0], nil
}

// Get a list of Log providers
func GetLogProviders(ctx *gin.Context, filters map[string]any) ([]*LogProvider, error) {
	// mapped `filterMapField, dbField`
	allowedFilters := [][]string{
		{"id", "id"},
		{"uuid", "uuid"},
		{"name", "name"},
	}
	tx := GetDBTransaction(ctx)
	providers := make([]*LogProvider, 0)
	tx = filter(filters, allowedFilters, tx)
	tx = tx.Find(&providers)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return providers, nil
}

// Get a single Alert Provider by UUID
func GetAlertProvider(ctx *gin.Context, uuid string) (*AlertProvider, error) {
	filters := make(map[string]any, 0)
	filters["uuid"] = uuid
	providers, err := GetAlertProviders(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(providers) == 0 {
		return nil, nil
	}
	return providers[0], nil
}

// Get a list of Alert providers
func GetAlertProviders(ctx *gin.Context, filters map[string]any) ([]*AlertProvider, error) {
	// mapped `filterMapField, dbField`
	allowedFilters := [][]string{
		{"id", "id"},
		{"uuid", "uuid"},
		{"name", "name"},
	}
	tx := GetDBTransaction(ctx)
	providers := make([]*AlertProvider, 0)
	tx = filter(filters, allowedFilters, tx)
	tx = tx.Find(&providers)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return providers, nil
}
