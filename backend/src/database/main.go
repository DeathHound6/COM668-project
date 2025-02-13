package database

import (
	"com668-backend/utility"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	conn             *gorm.DB    = nil
	defaultProviders []*Provider = []*Provider{
		{
			ID:   1,
			Name: "Sentry",
			Fields: []ProviderField{
				{
					ID:         1,
					ProviderID: 1,
					Key:        "enabled",
					Value:      "true",
					Type:       "bool",
					Required:   true,
				},
				{
					ID:         2,
					ProviderID: 1,
					Key:        "orgSlug",
					Value:      "testing-77",
					Type:       "string",
					Required:   true,
				},
				{
					ID:         3,
					ProviderID: 1,
					Key:        "projSlug",
					Value:      "test_app",
					Type:       "string",
					Required:   true,
				},
			},
			Type: "log",
		},
		{
			ID:   2,
			Name: "DynaTrace",
			Fields: []ProviderField{
				{
					ID:         4,
					ProviderID: 2,
					Key:        "enabled",
					Value:      "true",
					Type:       "bool",
					Required:   true,
				},
			},
			Type: "log",
		},
		{
			ID:   3,
			Name: "Slack",
			Fields: []ProviderField{
				{
					ID:         5,
					ProviderID: 3,
					Key:        "enabled",
					Value:      "true",
					Type:       "bool",
					Required:   true,
				},
			},
			Type: "alert",
		},
		{
			ID:   4,
			Name: "Microsoft Teams",
			Fields: []ProviderField{
				{
					ID:         6,
					ProviderID: 4,
					Key:        "enabled",
					Value:      "true",
					Type:       "bool",
					Required:   true,
				},
			},
			Type: "alert",
		},
	}
	defaultTeams []*Team = []*Team{
		{
			ID:   1,
			Name: "Engineering",
		},
		{
			ID:   2,
			Name: "DevOps",
		},
		{
			ID:   3,
			Name: "Monitoring",
		},
	}
	defaultUsers []*User = []*User{
		{
			ID:       1,
			Name:     "System",
			Email:    "test@example.com",
			Password: "system_user",
			Teams:    []Team{{ID: 2}},
			Admin:    true,
		},
		{
			ID:       2,
			Name:     "Test User",
			Email:    "user1@example.com",
			Password: "test_user",
			Teams:    []Team{{ID: 1}},
			Admin:    false,
		},
	}
	defaultHosts []*HostMachine = []*HostMachine{
		{
			ID:       1,
			Hostname: "test_app",
			IP4:      utility.Pointer("172.18.0.3"),
			IP6:      nil,
			OS:       "Linux",
			TeamID:   1,
		},
	}
	defaultIncidents []*Incident = []*Incident{
		{
			ID:              1,
			HostsAffected:   []HostMachine{{ID: 1}},
			Summary:         "Test Incident",
			Description:     "This is a test incident",
			Comments:        []IncidentComment{},
			CreatedAt:       time.Now(),
			ResolvedAt:      nil,
			ResolvedByID:    nil,
			ResolutionTeams: []Team{{ID: 2}},
		},
		{
			ID:            2,
			HostsAffected: []HostMachine{{ID: 1}},
			Summary:       "Test Incident 2",
			Description:   "This is a test incident",
			Comments: []IncidentComment{
				{
					Comment:       "This is a test comment",
					IncidentID:    2,
					CommentedByID: 1,
					CommentedAt:   time.Now().Add(time.Hour * -2),
				},
			},
			CreatedAt:       time.Now().Add(time.Hour * -3),
			ResolvedAt:      utility.Pointer(time.Now()),
			ResolvedByID:    utility.Pointer(uint(2)),
			ResolutionTeams: []Team{{ID: 2}},
		},
	}
)

func Connect() error {
	if conn != nil {
		return nil
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	params := map[string]any{
		"parseTime": "true",
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", username, password, host, dbName, strings.Join(utility.MapToSlice(params), "&"))
	log.Default().Println("Connecting to the database")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbl_",
			SingularTable: true,
			NoLowerCase:   false,
		},
	})
	if err != nil {
		return err
	}
	db = db.Session(&gorm.Session{Context: db.Statement.Context, NewDB: true})
	migrate(db)
	conn = db
	return nil
}

func GetDBConn() *gorm.DB {
	if conn == nil {
		Connect()
	}
	return conn
}

func migrate(conn *gorm.DB) {
	tx := conn.Begin()
	log.Default().Println("Migrating database")
	structs := []interface{}{
		Team{},
		User{},
		TeamUser{},
		Provider{},
		HostMachine{},
		Incident{},
		IncidentComment{},
		IncidentHost{},
		IncidentResolutionTeam{},
	}
	if gin.IsDebugging() {
		if err := tx.Migrator().DropTable(structs...); err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	if err := tx.AutoMigrate(structs...); err != nil {
		tx.Rollback()
		panic(err)
	}
	if gin.IsDebugging() {
		log.Default().Println("Inserting default data")
		if err := insert_default_data(tx); err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	tx = tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		panic(tx.Error)
	}
}

func insert_default_data(tx *gorm.DB) error {
	data := []interface{}{
		defaultTeams,
		defaultUsers,
		defaultProviders,
		defaultHosts,
		defaultIncidents,
	}
	for _, slice := range data {
		tx.Save(slice)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

func GetDBTransaction(ctx *gin.Context) *gorm.DB {
	tx, _ := ctx.Get("transaction")
	transaction := tx.(*gorm.DB)
	return transaction
}

func GetContext(tx *gorm.DB) *gin.Context {
	context, exists := tx.Get("context")
	if !exists {
		return nil
	}
	return context.(*gin.Context)
}

func handleError(ctx *gin.Context, err error) error {
	reqID := ctx.GetString("ReqID")
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("duplicate primary key was provided")
	} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("foreign key constraint violated")
	} else if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("an invalid enum value was given")
	} else {
		if mysqlErr, ok := err.(*mysql2.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				ctx.Set("errorCode", http.StatusBadRequest)
				return errors.New("duplicate field value was provided")
			default:
				log.Default().Printf("[%s] unhandled error: %e\n", reqID, err)
				ctx.Set("errorCode", http.StatusInternalServerError)
				return errors.New("an unhandled error occurred")
			}
		} else {
			log.Default().Printf("[%s] unhandled error: %e\n", reqID, err)
			ctx.Set("errorCode", http.StatusInternalServerError)
			return errors.New("an unhandled error occurred")
		}
	}
}
