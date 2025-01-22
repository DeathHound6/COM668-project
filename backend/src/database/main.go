package database

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
			Name:   "Sentry",
			Fields: "enabled;true;bool|orgSlug;testing-77;string|projSlug;test_app;string",
			Type:   "log",
		},
		{
			Name:   "Slack",
			Fields: "enabled;false;bool",
			Type:   "alert",
		},
	}
	defaultTeams []*Team = []*Team{
		{
			Name: "Engineering",
		},
	}
	defaultUsers []*User = []*User{
		{
			Name:     "System",
			Email:    "test@example.com",
			Password: "system_user",
			Teams: []Team{
				{
					Name: "Engineering",
				},
			},
			Admin: true,
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbName)
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
		Incident{},
		IncidentComment{},
	}
	if err := tx.Migrator().DropTable(structs...); err != nil {
		panic(err)
	}
	if err := tx.AutoMigrate(structs...); err != nil {
		panic(err)
	}
	log.Default().Println("Inserting default data")
	if err := insert_default_data(tx); err != nil {
		panic(err)
	}
	tx = tx.Commit()
	if tx.Error != nil {
		panic(tx.Error)
	}
}

func insert_default_data(tx *gorm.DB) error {
	data := []interface{}{
		defaultTeams,
		defaultUsers,
		defaultProviders,
	}
	for _, slice := range data {
		tx.Create(slice)
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
				log.Default().Printf("unhandled error: %e\n", err)
				ctx.Set("errorCode", http.StatusInternalServerError)
				return errors.New("an unhandled error occurred")
			}
		} else {
			log.Default().Printf("unhandled error: %e\n", err)
			ctx.Set("errorCode", http.StatusInternalServerError)
			return errors.New("an unhandled error occurred")
		}
	}
}

func filter(filters map[string]any, allowedFilters [][]string, tx *gorm.DB) *gorm.DB {
	for _, filterMap := range allowedFilters {
		value, ok := filters[filterMap[0]]
		if !ok {
			log.Default().Printf("Filter %s not allowed. skipping...\n", filterMap[0])
			continue
		}
		// Pagination filters
		if strings.ToLower(filterMap[0]) == "pagesize" {
			tx = tx.Limit(value.(int))
			continue
		}
		if strings.ToLower(filterMap[0]) == "page" {
			// Ensure page size exists - this allows us to calculate offset
			pageSize, ok := filters["pageSize"]
			if !ok {
				log.Default().Println("Page size not provided. skipping...")
				continue
			}
			// Value = page number
			page := value.(int) * pageSize.(int)
			tx = tx.Offset(page)
			continue
		}

		// Column filters
		tx = tx.Where(fmt.Sprintf("%s = ?", filterMap[1]), value)
	}
	return tx
}
