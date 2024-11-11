package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Incident struct {
	ID           uint       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID         string     `gorm:"column:uuid;size:36;unique;not null"`
	TeamID       uint       `gorm:"column:team_id;not null"`
	Team         Team       `gorm:"foreignKey:team_id;references:id"`
	Summary      string     `gorm:"column:summary;size:500;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime;not null"`
	ResolvedAt   *time.Time `gorm:"column:resolved_at"`
	ResolvedByID *uint      `gorm:"column:resolved_by_id"`
	ResolvedBy   *User      `gorm:"foreignKey:resolved_by_id;references:id"`
}

func (incident *Incident) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create an incident uuid")
	}
	if len(incident.Summary) > 500 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("summary cannot be greater than 500 characters")
	}
	incident.UUID = uuid
	return nil
}

func (incident *Incident) BeforeUpdate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	// If incident is being marked as resolved, set resolved time to now
	if incident.ResolvedBy != nil && incident.ResolvedAt == nil {
		time := time.Now()
		incident.ResolvedAt = &time
	}
	if len(incident.Summary) > 500 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("summary cannot be greater than 500 characters")
	}
	return nil
}

func (incident *Incident) BeforeDelete(tx *gorm.DB) error {
	ctx := GetContext(tx)
	ctx.Set("errorCode", http.StatusBadRequest)
	return errors.New("incidents cannot be deleted")
}

type IncidentComment struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	UUID          string    `gorm:"column:uuid;size:36;unique;not null"`
	Comment       string    `gorm:"column:comment;size:200;not null"`
	CommentedByID uint      `gorm:"column:commented_by_id;not null"`
	CommentedBy   User      `gorm:"foreignKey:commented_by_id;references:id"`
	CommentedAt   time.Time `gorm:"column:commented_at;autoCreateTime;not null"`
	IncidentID    uint      `gorm:"column:incident_id;not null"`
	Incident      Incident  `gorm:"foreignKey:incident_id;references:id"`
}

func (comment *IncidentComment) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create a comment uuid")
	}
	if len(comment.Comment) > 200 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("comment cannot be greater than 200 characters")
	}
	comment.UUID = uuid
	return nil
}

func (comment *IncidentComment) BeforeUpdate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if len(comment.Comment) > 200 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("comment cannot be greater than 200 characters")
	}
	return nil
}

func (comment *IncidentComment) BeforeDelete(tx *gorm.DB) error {
	return nil
}

func CreateIncident(ctx *gin.Context, body *utility.IncidentPostRequestBodySchema) (*Incident, error) {
	tx := GetDBTransaction(ctx)
	// todo: Get Team by id
	incident := &Incident{
		TeamID:  body.Team,
		Summary: body.Summary,
	}
	tx = tx.Create(incident)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return incident, nil
}
