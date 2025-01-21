package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	EmailRegexp string = "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+.[A-Za-z]{2,}"
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

func (user *User) hashPassword() (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	pass := string(bytes)
	return &pass, nil
}

func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create a user uuid")
	}
	if len(user.Name) > 30 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("user name cannot be greater than 30 characters")
	}
	if len(user.Email) > 30 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("user email cannot be greater than 30 characters")
	}
	if matched, err := regexp.MatchString(EmailRegexp, user.Email); !matched || err != nil {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("user email is not a valid email")
	}
	// NOTE: 72 chars is the max bcrypt supports
	if len(user.Password) > 72 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("user password cannot be greater than 72 characters")
	}
	password, err := user.hashPassword()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return err
	}
	if len(user.Teams) == 0 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("user must be part of at least 1 team")
	}
	user.Password = *password
	user.UUID = uuid
	return nil
}

type TeamUser struct {
	TeamID uint `gorm:"column:team_id;not null"`
	Team   Team `gorm:"foreignKey:team_id;references:id"`
	UserID uint `gorm:"column:user_id;not null"`
	User   User `gorm:"foreignKey:user_id;references:id"`
}

func GetUser(ctx *gin.Context, email string) (*User, error) {
	tx := GetDBTransaction(ctx)
	allowedFilters := [][]string{
		{"email", "email"},
	}
	filters := make(map[string]any, 0)
	filters["email"] = email
	tx = filter(filters, allowedFilters, tx)
	users := make([]*User, 0)
	// join to teams table
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
	tx := GetDBTransaction(ctx)
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
	tx = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(user)
	if tx.Error != nil {
		return nil, handleError(ctx, tx.Error)
	}
	return user, nil
}

func UpdateUser(ctx *gin.Context, user *User) error {
	tx := GetDBTransaction(ctx)
	tx = tx.Save(user)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

func DeleteUser(tx *gorm.DB) error {
	return nil
}
