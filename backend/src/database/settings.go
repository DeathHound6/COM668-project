package database

import (
	"strings"
)

// NOTE: DO NOT USE THIS IS A DB MODEL DIRECTLY
type provider struct {
	ID           uint              `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string            `gorm:"column:name"`
	ImageURL     string            `gorm:"column:image_url"`
	FieldsString string            `gorm:"column:fields"`
	Fields       map[string]string `gorm:"-"` // Tell gorm to ignore this field
}

// NOTE: DO NOT USE THIS IS A DB MODEL DIRECTLY
type providerSettings struct {
	ID             uint              `gorm:"column:id;primaryKey;autoIncrement"`
	ProviderID     uint              `gorm:"column:provider_id"`
	Provider       provider          `gorm:"foreignKey:provider_id;references:id;embedded;embeddedPrefix:provider_"`
	SettingsString string            `gorm:"column:settings"`
	Settings       map[string]string `gorm:"-"` // Tell gorm to ignore this field
}

// These are the actual db models - this is for gorm to identify table name
type LogProvider struct {
	provider
}
type LogProviderSettings struct {
	providerSettings
}
type AlertProvider struct {
	provider
}
type AlertProviderSettings struct {
	providerSettings
}

func (provider *provider) GetFields() {
	fields := make(map[string]string, 0)
	// Each field is separated by `|`
	dbFields := strings.Split(provider.FieldsString, "|")
	for _, field := range dbFields {
		// Each field is mapped `<key>=<value>`
		fieldKV := strings.Split(field, "=")
		fields[fieldKV[0]] = fieldKV[1]
	}
	provider.Fields = fields
}

func (settings *providerSettings) GetSettings() {
	fields := make(map[string]string, 0)
	// Each field is separated by `|`
	dbFields := strings.Split(settings.SettingsString, "|")
	for _, field := range dbFields {
		// Each field is mapped `<key>=<value>`
		fieldKV := strings.Split(field, "=")
		fields[fieldKV[0]] = fieldKV[1]
	}
	settings.Settings = fields
}
