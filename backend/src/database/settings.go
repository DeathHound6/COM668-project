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

type LogProviderSettings struct {
	ID         uint        `gorm:"column:id;primaryKey;autoIncrement"`
	ProviderID uint        `gorm:"column:provider_id;not null"`
	Provider   LogProvider `gorm:"foreignKey:provider_id;references:id"`
	Settings   string      `gorm:"column:settings;size:200;not null"`
}

type AlertProvider struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID   string `gorm:"column:uuid;size:36;unique;not null"`
	Name   string `gorm:"column:name;size:30;unique;not null"`
	Fields string `gorm:"column:fields;size:200;not null"`
}

type AlertProviderSettings struct {
	ID         uint          `gorm:"column:id;primaryKey;autoIncrement"`
	ProviderID uint          `gorm:"column:provider_id;not null"`
	Provider   AlertProvider `gorm:"foreignKey:provider_id;references:id"`
	Settings   string        `gorm:"column:settings;size:200;not null"`
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
	tx := GetDBTransaction(ctx)
	providers := make([]*LogProvider, 0)
	if id, ok := filters["id"]; ok {
		tx = tx.Where("id = ?", id)
	}
	if uuid, ok := filters["uuid"]; ok {
		tx = tx.Where("uuid = ?", uuid)
	}
	if name, ok := filters["name"]; ok {
		tx = tx.Where("name = ?", name)
	}
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
	tx := GetDBTransaction(ctx)
	providers := make([]*AlertProvider, 0)
	if id, ok := filters["id"]; ok {
		tx = tx.Where("id = ?", id)
	}
	if uuid, ok := filters["uuid"]; ok {
		tx = tx.Where("uuid = ?", uuid)
	}
	if name, ok := filters["name"]; ok {
		tx = tx.Where("name = ?", name)
	}
	out := tx.Find(&providers)
	if out.Error != nil {
		return nil, handleError(ctx, out.Error)
	}
	return providers, nil
}

func GetLogSettings(ctx *gin.Context, providerID uint) (*LogProviderSettings, error) {
	tx := GetDBTransaction(ctx)
	settings := make([]*LogProviderSettings, 0)
	out := tx.Where("provider_id=?", providerID).Find(&settings)
	if out.Error != nil {
		return nil, handleError(ctx, out.Error)
	}
	if len(settings) == 0 {
		return nil, nil
	}
	return settings[0], nil
}

func GetAlertSettings(ctx *gin.Context, providerID uint) (*AlertProviderSettings, error) {
	tx := GetDBTransaction(ctx)
	settings := make([]*AlertProviderSettings, 0)
	out := tx.Where("provider_id=?", providerID).Find(&settings)
	if out.Error != nil {
		return nil, handleError(ctx, out.Error)
	}
	if len(settings) == 0 {
		return nil, nil
	}
	return settings[0], nil
}
