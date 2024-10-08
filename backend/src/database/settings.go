package database

import (
	"github.com/gin-gonic/gin"
)

type LogProvider struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID     string `gorm:"column:uuid;size:16"`
	Name     string `gorm:"column:name;size:30"`
	ImageURL string `gorm:"column:image_url;size:100"`
	Fields   string `gorm:"column:fields;size:200"`
}

type LogProviderSettings struct {
	ID         uint        `gorm:"column:id;primaryKey;autoIncrement"`
	UUID       string      `gorm:"column:uuid;size:16"`
	ProviderID uint        `gorm:"column:provider_id"`
	Provider   LogProvider `gorm:"foreignKey:provider_id;references:id;embedded;embeddedPrefix:provider_"`
	Settings   string      `gorm:"column:settings;size:200"`
}

type AlertProvider struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID     string `gorm:"column:uuid;size:16"`
	Name     string `gorm:"column:name;size:30"`
	ImageURL string `gorm:"column:image_url;size:100"`
	Fields   string `gorm:"column:fields;size:200"`
}

type AlertProviderSettings struct {
	ID         uint          `gorm:"column:id;primaryKey;autoIncrement"`
	UUID       string        `gorm:"column:uuid;size:16"`
	ProviderID uint          `gorm:"column:provider_id"`
	Provider   AlertProvider `gorm:"foreignKey:provider_id;references:id;embedded;embeddedPrefix:provider_"`
	Settings   string        `gorm:"column:settings;size:200"`
}

func GetLogProviders(ctx *gin.Context) ([]*LogProvider, error) {
	logProviders := make([]*LogProvider, 0)

	return logProviders, nil
}
