package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Incident struct {
	ID         uint              `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Team       Team              `gorm:"foreignKey:id"`
	Summary    string            `json:"summary" gorm:"column:summary"`
	CreatedAt  time.Time         `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	ResolvedAt *time.Time        `json:"resolvedAt" gorm:"column:resolved_at"`
	ResolvedBy *User             `json:"resolvedBy" gorm:"column:resolved_by;foreignKey:id"`
	Comments   []IncidentComment `json:"comments" gorm:"foreignKey:id"`
}

func (incident *Incident) BeforeCreate(tx *gorm.DB) error {
	if len(incident.Summary) > 500 {
		return errors.New("summary cannot be greater than 500 characters")
	}
	return nil
}

func (incident *Incident) BeforeUpdate(tx *gorm.DB) error {
	// If incident is being marked as resolved, set resolved time to now
	if incident.ResolvedBy != nil {
		time := time.Now()
		incident.ResolvedAt = &time
	}
	return incident.BeforeCreate(tx)
}

func (incident *Incident) BeforeDelete(tx *gorm.DB) error {
	return errors.New("incidents cannot be deleted")
}

type IncidentComment struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Comment     string    `json:"comment" gorm:"column:comment"`
	CommentedBy User      `json:"commentedBy" gorm:"column:commented_by;foreignKey:id"`
	CommentedAt time.Time `json:"commentedAt" gorm:"column:commented_at"`
	Incident    Incident  `json:"incident" gorm:"foreignKey:id"`
}

func (comment *IncidentComment) BeforeCreate(tx *gorm.DB) error {
	if len(comment.Comment) > 200 {
		return errors.New("comment cannot be greater than 200 characters")
	}
	return nil
}

func (comment *IncidentComment) BeforeUpdate(tx *gorm.DB) error {
	return comment.BeforeCreate(tx)
}

func (comment *IncidentComment) BeforeDelete(tx *gorm.DB) error {
	return nil
}
