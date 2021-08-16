package database

import (
	"fmt"
	"log"

	"github.com/energy-uktc/grouping-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//PostgresqlConfiguration which should be readed from the config.yaml
type PostgresqlConfiguration struct {
	host         string `mapstructure:"host"`
	port         int    `mapstructure:"port"`
	databaseName string `mapstructure:"databaseName"`
	user         string `mapstructure:"user"`
	password     string `mapstructure:"password"`
	sslmode      string `mapstructure:"sslmode"`
	timeZone     string `mapstructure:"timeZone"`
}

func (c *PostgresqlConfiguration) LoadConfiguration() {
	c.host = readPostgreValueFromConfig("host", true).(string)
	c.port = readPostgreValueFromConfig("port", true).(int)
	c.databaseName = readPostgreValueFromConfig("databaseName", true).(string)
	c.user = readPostgreValueFromConfig("user", true).(string)
	c.password = readPostgreValueFromConfig("password", true).(string)
	c.sslmode = readPostgreValueFromConfig("sslMode", false).(string)
	c.timeZone = readPostgreValueFromConfig("timeZone", false).(string)

}

func readPostgreValueFromConfig(key string, important bool) interface{} {
	value := config.GetPath("postgresql." + key)
	if value == nil && important {
		log.Fatalf("Postgresql %s is not configured", key)
	}
	return value
}

//CreateConnection returns posgres dialector with connection to the database
func (c *PostgresqlConfiguration) CreateConnection() gorm.Dialector {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d ", c.host, c.user, c.password, c.databaseName, c.port)
	if c.sslmode != "" {
		connectionString += "sslmode=" + c.sslmode + " "
	}
	if c.timeZone != "" {
		connectionString += "TimeZone=" + c.timeZone + " "
	}
	return postgres.Open(connectionString)
}
