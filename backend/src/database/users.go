package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID     string `gorm:"column:uuid;size:36;unique;not null"`
	Name     string `gorm:"column:name;size:30;unique;not null"`
	Email    string `gorm:"column:email;size:30;unique;not null"`
	Password string `gorm:"column:password;size:72;not null"`
	Admin    bool   `gorm:"column:admin;not null"`
	Teams    []Team `gorm:"many2many:team_user"`
	SlackID  string `gorm:"column:slack_id;size:20"`
}

func (user *User) hashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if user.UUID == "" {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusInternalServerError)
			}
			return errors.New("failed to create a user uuid")
		}
		user.UUID = uuid
	}
	password, err := user.hashPassword()
	if err != nil {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusInternalServerError)
		}
		return err
	}
	user.Password = password
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	password, err := user.hashPassword()
	if err != nil {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusInternalServerError)
		}
		return err
	}
	user.Password = password
	return nil
}

type TeamUser struct {
	TeamID uint `gorm:"column:team_id;primaryKey"`
	Team   Team `gorm:"foreignKey:team_id;references:id"`
	UserID uint `gorm:"column:user_id;primaryKey"`
	User   User `gorm:"foreignKey:user_id;references:id"`
}

type GetUserFilters struct {
	UUID  *string
	Email *string
}

func GetUser(ctx *gin.Context, filters GetUserFilters) (*User, error) {
	tx := GetDBTransaction(ctx).Model(&User{})
	if filters.UUID != nil {
		tx = tx.Where("uuid = ?", *filters.UUID)
	}
	if filters.Email != nil {
		tx = tx.Where("email = ?", *filters.Email)
	}
	users := make([]*User, 0)
	tx = tx.Preload("Teams")
	tx = tx.Find(&users)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func CreateUser(ctx *gin.Context, body *utility.UserPostRequestBodySchema) (*User, error) {
	tx := GetDBTransaction(ctx).Model(&User{})
	teams := make([]Team, 0)
	for _, team := range body.Teams {
		teams = append(teams, Team{
			UUID: team,
		})
	}
	user := &User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Teams:    teams,
	}
	tx = tx.Clauses(clause.OnConflict{DoNothing: true}).Save(user)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return user, nil
}

func UpdateUser(ctx *gin.Context, user *User) error {
	tx := GetDBTransaction(ctx).Model(&User{})
	tx = tx.Save(user)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

func DeleteUser(tx *gorm.DB) error {
	return nil
}
