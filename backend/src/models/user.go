package models

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	EmailRegexp string = "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+.[A-Za-z]{2,}"
)

type User struct {
	ID       uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"column:name;size:30"`
	Email    string `json:"email" gorm:"column:email;size:30"`
	Password string `json:"password" gorm:"column:password;size:72"`
	Teams    []Team `json:"teams" gorm:"foreignKey:id"`
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
	ID    uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"column:name"`
	Users []User `json:"users" gorm:"foreignKey:id"`
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
