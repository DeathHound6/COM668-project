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
	UUID            string            `gorm:"column:uuid;size:36;unique;not null;uniqueIndex"`
	HostsAffected   []HostMachine     `gorm:"many2many:incident_host"`
	Description     string            `gorm:"column:description;size:500"`
	Summary         string            `gorm:"column:summary;size:100;not null"`
	Comments        []IncidentComment `gorm:"foreignKey:incident_id;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time         `gorm:"column:created_at;autoCreateTime;not null"`
	ResolvedAt      *time.Time        `gorm:"column:resolved_at;index"`
	ResolvedByID    *uint             `gorm:"column:resolved_by_id"`
	ResolvedBy      *User             `gorm:"foreignKey:resolved_by_id;references:id"`
	ResolutionTeams []Team            `gorm:"many2many:incident_resolution_team"`
	Hash            string            `gorm:"column:hash;size:64;not null;uniqueIndex"`
}

func (incident *Incident) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if incident.UUID == "" {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusInternalServerError)
			}
			return errors.New("failed to create an incident uuid")
		}
		incident.UUID = uuid
	}
	return nil
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
	if comment.UUID == "" {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusInternalServerError)
			}
			return errors.New("failed to create a comment uuid")
		}
		comment.UUID = uuid
	}
	return nil
}

type IncidentHost struct {
	ID            uint        `gorm:"column:id;primaryKey;autoIncrement"`
	IncidentID    uint        `gorm:"column:incident_id"`
	Incident      Incident    `gorm:"foreignKey:incident_id;references:id"`
	HostMachineID uint        `gorm:"column:host_machine_id"`
	HostMachine   HostMachine `gorm:"foreignKey:host_machine_id;references:id"`
}
type IncidentResolutionTeam struct {
	ID         uint     `gorm:"column:id;primaryKey;autoIncrement"`
	IncidentID uint     `gorm:"column:incident_id"`
	Incident   Incident `gorm:"foreignKey:incident_id;references:id"`
	TeamID     uint     `gorm:"column:team_id"`
	Team       Team     `gorm:"foreignKey:team_id;references:id"`
}

type GetIncidentsFilters struct {
	MyTeams  bool
	UUID     *string
	Resolved *bool
	Page     *int
	PageSize *int
	Hash     *string
}

func GetIncident(ctx *gin.Context, uuid string) (*Incident, error) {
	incidents, count, err := GetIncidents(ctx, GetIncidentsFilters{
		PageSize: utility.Pointer(1),
		MyTeams:  false,
		UUID:     &uuid,
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		ctx.Set("errorCode", http.StatusNotFound)
		return nil, errors.New("incident not found")
	}
	return incidents[0], nil
}

func GetIncidents(ctx *gin.Context, filters GetIncidentsFilters) ([]*Incident, int64, error) {
	tx := GetDBTransaction(ctx).Model(&Incident{})
	tx = tx.Preload("HostsAffected").
		Preload("HostsAffected.Team").
		Preload("ResolutionTeams").
		Preload("ResolutionTeams.Users").
		Preload("ResolvedBy").
		Preload("ResolvedBy.Teams").
		Preload("Comments", func(t *gorm.DB) *gorm.DB {
			// get comments in descending order by created_at timestamp
			return t.Order("commented_at DESC")
		}).
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
	if filters.MyTeams {
		user := ctx.MustGet("user").(*User)
		tx = tx.Joins("LEFT JOIN tbl_incident_resolution_team ON tbl_incident_resolution_team.incident_id = tbl_incident.id").
			Joins("LEFT JOIN tbl_incident_host ON tbl_incident_host.incident_id = tbl_incident.id").
			Joins("LEFT JOIN tbl_host_machine ON tbl_host_machine.id = tbl_incident_host.host_machine_id").
			Joins("LEFT JOIN tbl_team_user ON tbl_team_user.team_id = tbl_incident_resolution_team.team_id OR tbl_team_user.team_id = tbl_host_machine.team_id").
			Joins("LEFT JOIN tbl_user ON tbl_user.id = tbl_team_user.user_id").
			Where("tbl_user.uuid = ?", user.UUID)
	}
	if filters.UUID != nil {
		tx = tx.Where("tbl_incident.uuid = ?", *filters.UUID)
	}
	if filters.Hash != nil {
		tx = tx.Where("tbl_incident.hash = ?", *filters.Hash)
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
	incident := &Incident{
		Summary:     body.Summary,
		Description: body.Description,
		CreatedAt:   time.Now(),
		Hash:        body.Hash,
	}
	tx = tx.Create(incident)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}

	// insert hosts - m2m
	hosts := make([]*IncidentHost, 0)
	hs, count, err := GetHosts(ctx, GetHostsFilters{
		UUIDs:    body.HostsAffected,
		PageSize: utility.Pointer(len(body.HostsAffected)),
	})
	if err != nil {
		return nil, err
	}
	if int(count) != len(body.HostsAffected) {
		ctx.Set("errorCode", http.StatusBadRequest)
		return nil, errors.New("one or more hosts not found")
	}
	for _, host := range hs {
		hosts = append(hosts, &IncidentHost{
			IncidentID:    incident.ID,
			HostMachineID: host.ID,
		})
	}
	tx = GetDBTransaction(ctx).Model(&IncidentHost{}).CreateInBatches(hosts, 1)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}

	// insert resolution teams - m2m
	teams := make([]*IncidentResolutionTeam, 0)
	ts, count, err := GetTeams(ctx, GetTeamsFilters{
		UUIDs:    body.ResolutionTeams,
		PageSize: utility.Pointer(len(body.ResolutionTeams)),
	})
	if err != nil {
		return nil, err
	}
	if int(count) != len(body.ResolutionTeams) {
		ctx.Set("errorCode", http.StatusBadRequest)
		return nil, errors.New("one or more teams not found")
	}
	for _, team := range ts {
		teams = append(teams, &IncidentResolutionTeam{
			IncidentID: incident.ID,
			TeamID:     team.ID,
		})
	}
	tx = GetDBTransaction(ctx).Model(&IncidentResolutionTeam{}).CreateInBatches(teams, 1)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return incident, nil
}

func UpdateIncident(ctx *gin.Context, filters GetIncidentsFilters, incident *Incident) error {
	tx := GetDBTransaction(ctx)

	temp1 := tx.Model(&Incident{})
	if filters.UUID != nil {
		temp1 = temp1.Where("uuid = ?", *filters.UUID)
	}
	fields := map[string]any{
		"summary":        incident.Summary,
		"description":    incident.Description,
		"resolved_at":    incident.ResolvedAt,
		"resolved_by_id": incident.ResolvedByID,
	}
	if err := temp1.Updates(fields).Error; err != nil {
		return handleError(ctx, err)
	}

	// replace hosts - m2m
	hosts := make([]IncidentHost, 0)
	for _, host := range incident.HostsAffected {
		hosts = append(hosts, IncidentHost{
			IncidentID:    incident.ID,
			HostMachineID: host.ID,
		})
	}
	if err := tx.Model(&IncidentHost{}).Where("incident_id = ?", incident.ID).Delete(&IncidentHost{}).Error; err != nil {
		return handleError(ctx, err)
	}
	if len(hosts) > 0 {
		if err := tx.Model(&IncidentHost{}).CreateInBatches(hosts, 1).Error; err != nil {
			return handleError(ctx, err)
		}
	}

	// replace resolution teams - m2m
	teams := make([]IncidentResolutionTeam, 0)
	for _, team := range incident.ResolutionTeams {
		teams = append(teams, IncidentResolutionTeam{
			IncidentID: incident.ID,
			TeamID:     team.ID,
		})
	}
	if err := tx.Model(&IncidentResolutionTeam{}).Where("incident_id = ?", incident.ID).Delete(&IncidentResolutionTeam{}).Error; err != nil {
		return handleError(ctx, err)
	}
	if len(teams) > 0 {
		tx = tx.CreateInBatches(teams, 1)
		if tx.Error != nil {
			return handleError(ctx, tx.Error)
		}
	}
	return nil
}

func GetIncidentComment(ctx *gin.Context, incidentUUID, commentUUID string) (*IncidentComment, error) {
	incident, err := GetIncident(ctx, incidentUUID)
	if err != nil {
		return nil, err
	}
	var comment *IncidentComment = nil
	for _, c := range incident.Comments {
		if c.UUID == commentUUID {
			comment = &c
			break
		}
	}
	if comment == nil {
		ctx.Set("errorCode", http.StatusNotFound)
		return nil, errors.New("comment not found")
	}
	return comment, nil
}

func CreateIncidentComment(ctx *gin.Context, comment *IncidentComment) (*IncidentComment, error) {
	tx := GetDBTransaction(ctx).Model(&IncidentComment{})
	tx = tx.Create(comment)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return comment, nil
}

func DeleteIncidentComment(ctx *gin.Context, uuid string) error {
	tx := GetDBTransaction(ctx).Model(&IncidentComment{})
	tx = tx.Where("uuid = ?", uuid)
	tx = tx.Delete(&IncidentComment{})
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}
