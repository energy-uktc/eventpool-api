package database

import (
	"log"
	"time"

	"github.com/energy-uktc/eventpool-api/config"
	"gorm.io/gorm"
)

//Configuration interface for database configuration
type Configuration interface {
	CreateConnection() gorm.Dialector
}

//DbConn Connection to the database
var DbConn *gorm.DB

//SetupDatabase setups connection to the db
func SetupDatabase() {
	var err error
	dbConfig := getDatabaseConfiguration()
	DbConn, err = gorm.Open(dbConfig.CreateConnection(), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDb, err := DbConn.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(60 * time.Second)
	DbConn.Set("gorm:auto_preload", true)
}

func getDatabaseConfiguration() Configuration {
	if config.Properties.Database == "postgresql" {
		postgresql := &PostgresqlConfiguration{}
		postgresql.LoadConfiguration()
		return postgresql
	} else if config.Properties.Database == "cloudpostgresql" {
		cloudDb := &CloudSqlConfiguration{}
		cloudDb.LoadConfiguration()
		return cloudDb
	}
	panic("No database configuration found")
}
