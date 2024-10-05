package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Incident struct {
	ID           uint       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID         string     `gorm:"column:uuid"`
	TeamID       uint       `gorm:"column:team_id"`
	Team         Team       `gorm:"foreignKey:team_id;references:id"`
	Summary      string     `gorm:"column:summary"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	ResolvedAt   *time.Time `gorm:"column:resolved_at"`
	ResolvedByID uint       `gorm:"column:resolved_by_id"`
	ResolvedBy   *User      `gorm:"foreignKey:resolved_by_id;references:id"`
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
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Comment       string    `gorm:"column:comment"`
	CommentedByID uint      `gorm:"column:commented_by_id"`
	CommentedBy   User      `gorm:"foreignKey:commented_by_id;references:id"`
	CommentedAt   time.Time `gorm:"column:commented_at;autoCreateTime"`
	IncidentID    uint      `gorm:"column:incident_id"`
	Incident      Incident  `gorm:"foreignKey:incident_id;references:id"`
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
