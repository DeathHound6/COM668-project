package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	conn *gorm.DB = nil
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
	log.Default().Println("Migrating database")
	structs := []interface{}{
		Team{},
		User{},
		TeamUser{},
		LogProvider{},
		LogProviderSettings{},
		AlertProvider{},
		AlertProviderSettings{},
		Incident{},
		IncidentComment{},
	}
	if err := conn.AutoMigrate(structs...); err != nil {
		panic(err)
	}
	conn.Commit()
}

func GetDBResults(tx *gorm.DB, model interface{}) ([]interface{}, error) {
	rows, err := tx.Rows()
	results := make([]interface{}, 0)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(model); err != nil {
			return nil, err
		}
		results = append(results, model)
	}
	return results, nil
}
