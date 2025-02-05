package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Team struct {
	ID    uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID  string `gorm:"column:uuid;size:36;unique;not null"`
	Name  string `gorm:"column:name;size:30;unique;not null"`
	Users []User `gorm:"many2many:team_user"`
}

func (team *Team) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create an incident uuid")
	}
	if len(team.Name) > 30 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("team name cannot be greater than 30 characters")
	}
	team.UUID = uuid
	return nil
}

func (team *Team) BeforeUpdate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if len(team.Name) > 30 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("team name cannot be greater than 30 characters")
	}
	return nil
}

func (team *Team) BeforeDelete(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if len(team.Users) > 0 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("teams cannot be deleted if there are still users in them")
	}
	return nil
}

type GetTeamsFilters struct {
	UUID     *string
	Page     *int
	PageSize *int
}

func GetTeam(ctx *gin.Context, filters GetTeamsFilters) (*Team, error) {
	teams, count, err := GetTeams(ctx, GetTeamsFilters{
		UUID: filters.UUID,
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		ctx.Set("errorCode", http.StatusNotFound)
		return nil, errors.New("team not found")
	}
	return teams[0], nil
}

func GetTeams(ctx *gin.Context, filters GetTeamsFilters) ([]*Team, int64, error) {
	tx := GetDBTransaction(ctx).Model(&Team{})

	if filters.UUID != nil {
		tx = tx.Where("uuid = ?", *filters.UUID)
	}

	var count int64
	tx = tx.Count(&count)
	if filters.PageSize != nil {
		tx = tx.Limit(*filters.PageSize)
		if filters.Page != nil {
			tx = tx.Offset((*filters.Page - 1) * *filters.PageSize)
		}
	}

	var teams []*Team
	tx = tx.Find(&teams)
	if tx.Error != nil {
		return nil, -1, handleError(ctx, tx.Error)
	}
	return teams, count, nil
}

func CreateTeam(ctx *gin.Context, body *utility.TeamPostRequestBodySchema) error {
	return nil
}
