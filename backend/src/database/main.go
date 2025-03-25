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
			UUID: "0a846a37-d039-42c6-a1c9-699763ae646e",
			Name: "Sentry",
			Type: "log",
		},
		{
			UUID: "31a3e142-1222-45ca-9c89-91c736cdf4a6",
			Name: "DynaTrace",
			Type: "log",
		},
		{
			UUID: "24d0e277-1f85-4779-9d25-3db24055b493",
			Name: "Slack",
			Type: "alert",
		},
		{
			UUID: "87c6c563-21df-4a3e-8746-f7871e6fa431",
			Name: "Microsoft Teams",
			Type: "alert",
		},
	}
	defaultProvidersFields []*ProviderField = []*ProviderField{
		{
			ProviderID: 1,
			Key:        "enabled",
			Value:      "true",
			Type:       "bool",
			Required:   true,
		},
		{
			ProviderID: 1,
			Key:        "orgSlug",
			Value:      "testing-77",
			Type:       "string",
			Required:   true,
		},
		{
			ProviderID: 1,
			Key:        "projSlug",
			Value:      "test_app",
			Type:       "string",
			Required:   true,
		},
		{
			ProviderID: 2,
			Key:        "enabled",
			Value:      "true",
			Type:       "bool",
			Required:   true,
		},
		{
			ProviderID: 3,
			Key:        "enabled",
			Value:      "true",
			Type:       "bool",
			Required:   true,
		},
		{
			ProviderID: 4,
			Key:        "enabled",
			Value:      "true",
			Type:       "bool",
			Required:   true,
		},
	}
	defaultTeams []*Team = []*Team{
		{
			UUID: "574b5d6a-1fcd-43bf-bb31-7e870ca458d4",
			Name: "App 1",
		},
		{
			UUID: "7544efed-9e6f-4bf1-a8f9-9a93df5944df",
			Name: "DevOps",
		},
		{
			UUID: "89e7fdc7-dd8c-471c-9e23-94cf678412a2",
			Name: "NetOps",
		},
	}
	defaultUsers []*User = []*User{
		{
			UUID:     "39ab8bf8-fd8c-43c2-b691-3acb4f5a3fab",
			Name:     "System",
			Email:    "test@example.com",
			Password: "system_user",
			Admin:    true,
		},
		{
			UUID:     "417cd42a-ddff-42dc-b358-801807522dbd",
			Name:     "Test User",
			Email:    "user1@example.com",
			Password: "test_user",
			Admin:    false,
		},
	}
	defaultTeamUsers []*TeamUser = []*TeamUser{
		{
			UserID: 1,
			TeamID: 2,
		},
		{
			UserID: 2,
			TeamID: 1,
		},
	}
	defaultHosts []*HostMachine = []*HostMachine{
		{
			UUID:     "c3cb5381-7b79-4bbe-9337-8a27f94646a4",
			Hostname: "7e83c1b6c515",
			IP4:      utility.Pointer("172.18.0.3"),
			IP6:      nil,
			OS:       "Linux",
			TeamID:   1,
		},
	}
	defaultIncidents []*Incident = []*Incident{
		{
			UUID:         "eddb82c0-50fd-4faa-adc0-95db00df52da",
			Summary:      "Test Incident",
			Description:  "This is a test incident",
			Comments:     []IncidentComment{},
			CreatedAt:    time.Now(),
			ResolvedAt:   nil,
			ResolvedByID: nil,
			Hash:         "b10c4cf0a834c7f5b07f395c633d906fe7d969b9",
		},
		{
			UUID:         "30daaadd-596c-4676-9194-c8f48a654931",
			Summary:      "Test Incident 2",
			Description:  "This is a test incident",
			CreatedAt:    time.Now().Add(time.Hour * -3),
			ResolvedAt:   utility.Pointer(time.Now()),
			ResolvedByID: utility.Pointer(uint(2)),
			Hash:         "5326ca01ce201c73f9b5d112c7d4c6eb0d14abbc",
		},
	}
	defaultIncidentComments []*IncidentComment = []*IncidentComment{
		{
			UUID:          "3ca620e9-b90f-492d-9597-5677f762ec00",
			Comment:       "This is a test comment",
			IncidentID:    2,
			CommentedByID: 1,
			CommentedAt:   time.Now().Add(time.Hour * -2),
		},
	}
	defaultIncidentResolutionTeams []*IncidentResolutionTeam = []*IncidentResolutionTeam{
		{
			IncidentID: 1,
			TeamID:     2,
		},
		{
			IncidentID: 2,
			TeamID:     2,
		},
	}
	defaultIncidentHosts []*IncidentHost = []*IncidentHost{
		{
			IncidentID:    1,
			HostMachineID: 1,
		},
		{
			IncidentID:    2,
			HostMachineID: 1,
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
	structs := []interface{}{
		Team{},
		User{},
		TeamUser{},
		Provider{},
		ProviderField{},
		HostMachine{},
		Incident{},
		IncidentComment{},
		IncidentHost{},
		IncidentResolutionTeam{},
	}
	if gin.IsDebugging() {
		log.Default().Println("Dropping tables")
		if err := tx.Migrator().DropTable(structs...); err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	log.Default().Println("Migrating database")
	if err := tx.AutoMigrate(structs...); err != nil {
		tx.Rollback()
		panic(err)
	}
	if gin.IsDebugging() {
		log.Default().Println("Inserting default data")
		if err := insertDefaultData(tx); err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	err := tx.SetupJoinTable(&Incident{}, "HostsAffected", &IncidentHost{})
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	err = tx.SetupJoinTable(&Incident{}, "ResolutionTeams", &IncidentResolutionTeam{})
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	tx = tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		panic(tx.Error)
	}
}

func insertDefaultData(tx *gorm.DB) error {
	data := []any{
		defaultTeams,
		defaultUsers,
		defaultTeamUsers,
		defaultProviders,
		defaultProvidersFields,
		defaultHosts,
		defaultIncidents,
		defaultIncidentComments,
		defaultIncidentResolutionTeams,
		defaultIncidentHosts,
	}
	for _, slice := range data {
		log.Default().Printf("Inserting records %T\n", slice)
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
