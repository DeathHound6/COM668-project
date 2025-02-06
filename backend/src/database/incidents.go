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
	ID              uint              `gorm:"column:id;primaryKey;autoIncrement"`
	UUID            string            `gorm:"column:uuid;size:36;unique;not null"`
	HostsAffected   []HostMachine     `gorm:"many2many:incident_host"`
	Summary         string            `gorm:"column:summary;size:500;not null"`
	Comments        []IncidentComment `gorm:"foreignKey:incident_id"`
	CreatedAt       time.Time         `gorm:"column:created_at;autoCreateTime;not null"`
	ResolvedAt      *time.Time        `gorm:"column:resolved_at"`
	ResolvedByID    *uint             `gorm:"column:resolved_by_id"`
	ResolvedBy      *User             `gorm:"foreignKey:resolved_by_id;references:id"`
	ResolutionTeams []Team            `gorm:"many2many:incident_resolution_team"`
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

type IncidentHost struct {
	IncidentID    uint        `gorm:"column:incident_id;primaryKey"`
	Incident      Incident    `gorm:"foreignKey:incident_id;references:id"`
	HostMachineID uint        `gorm:"column:host_machine_id;primaryKey"`
	HostMachine   HostMachine `gorm:"foreignKey:host_machine_id;references:id"`
}

type GetIncidentsFilters struct {
	Resolved *bool
	Page     *int
	PageSize *int
}

func GetIncidents(ctx *gin.Context, filters GetIncidentsFilters) ([]*Incident, int64, error) {
	tx := GetDBTransaction(ctx).Model(&Incident{})
	tx = tx.Preload("HostsAffected").
		Preload("HostsAffected.Team").
		Preload("ResolutionTeams").
		Preload("ResolvedBy").
		Preload("ResolvedBy.Teams").
		Preload("Comments").
		Preload("Comments.CommentedBy").
		Preload("Comments.CommentedBy.Teams")
	incidents := make([]*Incident, 0)

	// apply filters
	if filters.Resolved != nil {
		if *filters.Resolved {
			tx = tx.Where("resolved_at IS NOT NULL")
		} else {
			tx = tx.Where("resolved_at IS NULL")
		}
	}

	var count int64
	tx.Count(&count)
	if filters.PageSize != nil {
		tx = tx.Limit(*filters.PageSize)
		if filters.Page != nil {
			tx = tx.Offset(*filters.PageSize * (*filters.Page - 1))
		}
	}

	tx = tx.Find(&incidents)
	if tx.Error != nil {
		return nil, -1, handleError(ctx, tx.Error)
	}
	return incidents, count, nil
}

func CreateIncident(ctx *gin.Context, body *utility.IncidentPostRequestBodySchema) (*Incident, error) {
	tx := GetDBTransaction(ctx).Model(&Incident{})
	hosts := make([]HostMachine, 0)
	// todo: validate hosts
	for _, host := range body.HostsAffected {
		hosts = append(hosts, HostMachine{
			UUID: host,
		})
	}
	incident := &Incident{
		HostsAffected: hosts,
		Summary:       body.Summary,
	}
	tx = tx.Create(incident)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return incident, nil
}
