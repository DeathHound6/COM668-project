package database

import (
	"com668-backend/utility"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	EmailRegexp string = "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+.[A-Za-z]{2,}"
)

type User struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name     string `gorm:"column:name;size:30"`
	Email    string `gorm:"column:email;size:30"`
	Password string `gorm:"column:password;size:72"`
	Teams    []Team `gorm:"many2many:team_user"`
}

func (user *User) hashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if len(user.Name) > 30 {
		return errors.New("user name cannot be greater than 30 characters")
	}
	if len(user.Email) > 30 {
		return errors.New("user email cannot be greater than 30 characters")
	}
	if matched, err := regexp.MatchString(EmailRegexp, user.Email); !matched || err != nil {
		if err != nil {
			return err
		}
		return errors.New("user email is not a valid email")
	}
	// NOTE: 72 chars is the max bcrypt supports
	if len(user.Password) > 72 {
		return errors.New("user password cannot be greater than 72 characters")
	}
	if err := user.hashPassword(); err != nil {
		return err
	}
	if len(user.Teams) == 0 {
		return errors.New("user must be part of at least 1 team")
	}
	return nil
}

type Team struct {
	ID    uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name  string `gorm:"column:name"`
	Users []User `gorm:"many2many:team_user"`
}

func (team *Team) BeforeCreate(tx *gorm.DB) error {
	if len(team.Name) > 30 {
		return errors.New("team name cannot be greater than 30 characters")
	}
	return nil
}

func (team *Team) BeforeUpdate(tx *gorm.DB) error {
	return team.BeforeCreate(tx)
}

func (team *Team) BeforeDelete(tx *gorm.DB) error {
	if len(team.Users) > 0 {
		return errors.New("teams cannot be deleted if there are still users in them")
	}
	return nil
}

type TeamUser struct {
	TeamID uint `gorm:"column:team_id"`
	Team   Team `gorm:"foreignKey:team_id;references:id"`
	UserID uint `gorm:"column:user_id"`
	User   User `gorm:"foreignKey:user_id;references:id"`
}

func GetUser(tx *gorm.DB) error {
	return nil
}

func CreateUser(tx *gorm.DB, body *utility.UserPostRequestBodySchema) error {
	teams := make([]Team, 0)
	for _, teamName := range body.Teams {
		teams = append(teams, Team{
			Name: teamName,
		})
	}
	user := &User{
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
