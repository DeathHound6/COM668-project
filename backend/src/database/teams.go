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
	UUID  string `gorm:"column:uuid;size:16"`
	Name  string `gorm:"column:name"`
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

func CreateTeam(ctx *gin.Context, body *utility.TeamPostRequestBodySchema) error {
	return nil
}

func getTeams(tx *gorm.DB, teamNames []string) []Team {
	teams := make([]Team, 0)
	tx.Where("name IN (?)", teamNames).Take(&teams)
	return teams
}
