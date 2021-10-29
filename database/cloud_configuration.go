package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CloudSqlConfiguration struct {
	instanceConnectionName string
	port                   int
	databaseName           string
	user                   string
	password               string
}

func (c *CloudSqlConfiguration) LoadConfiguration() {
	c.instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
	c.port = readPostgreValueFromConfig("port", true).(int)
	c.databaseName = mustGetenv("DB_NAME")
	c.user = mustGetenv("DB_USER")
	c.password = mustGetenv("DB_PASS")
}

//CreateConnection returns posgres dialector with connection to the database
func (c *CloudSqlConfiguration) CreateConnection() gorm.Dialector {

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	connectionString := fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", c.user, c.password, c.databaseName, socketDir, c.instanceConnectionName)
	return postgres.Open(connectionString)
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}
