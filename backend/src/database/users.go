package database

import (
	"com668-backend/models"
	"com668-backend/utility"
	"errors"

	"gorm.io/gorm"
)

func GetUser(tx *gorm.DB) error {
	return nil
}

func CreateUser(tx *gorm.DB, body *utility.UserPostRequestBodySchema) error {
	teams := make([]models.Team, 0)
	for _, teamName := range body.Teams {
		teams = append(teams, models.Team{
			Name: teamName,
		})
	}
	user := &models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Teams:    teams,
	}
	out := tx.Create(user)
	if out.Error != nil || out.RowsAffected == 0 {
		if out.Error != nil {
			return out.Error
		}
		return errors.New("failed to create new user")
	}
	return nil
}

func UpdateUser(tx *gorm.DB) error {
	return nil
}

func DeleteUser(tx *gorm.DB) error {
	return nil
}
