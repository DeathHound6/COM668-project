package models

import (
	"strings"
)

// NOTE: DO NOT USE THIS IS A DB MODEL DIRECTLY
type Provider struct {
	ID           uint              `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name         string            `json:"name" gorm:"column:name"`
	ImageURL     string            `json:"imageUrl" gorm:"column:image_url"`
	FieldsString string            `gorm:"column:fields"`
	Fields       map[string]string `json:"fields"`
}

// NOTE: DO NOT USE THIS IS A DB MODEL DIRECTLY
type ProviderSettings struct {
	ID             uint              `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Provider       Provider          `gorm:"foreignKey:id"`
	SettingsString string            `gorm:"column:settings"`
	Settings       map[string]string `json:"settings"`
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
