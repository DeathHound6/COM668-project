package models

import (
	"strings"
)

// NOTE: DO NOT USE THIS IS A DB MODEL DIRECTLY
type Provider struct {
	ID           uint              `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string            `gorm:"column:name"`
	ImageURL     string            `gorm:"column:image_url"`
	FieldsString string            `gorm:"column:fields"`
	Fields       map[string]string `gorm:"-"` // Tell gorm to ignore this field
}

// NOTE: DO NOT USE THIS IS A DB MODEL DIRECTLY
type ProviderSettings struct {
	ID             uint              `gorm:"column:id;primaryKey;autoIncrement"`
	ProviderID     uint              `gorm:"column:provider_id"`
	Provider       Provider          `gorm:"foreignKey:provider_id;references:id;embedded;embeddedPrefix:provider_"`
	SettingsString string            `gorm:"column:settings"`
	Settings       map[string]string `gorm:"-"` // Tell gorm to ignore this field
}

// These are the actual db models - this is for gorm to identify table name
type LogProvider struct {
	Provider
}
type LogProviderSettings struct {
	ProviderSettings
}
type AlertProvider struct {
	Provider
}
type AlertProviderSettings struct {
	ProviderSettings
}

func (provider *Provider) GetFields() {
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

func (settings *ProviderSettings) GetSettings() {
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
