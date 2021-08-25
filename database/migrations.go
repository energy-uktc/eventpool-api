package database

import (
	"fmt"

	"github.com/energy-uktc/eventpool-api/entities"
)

//UpdateSchema performs database migrations for used entities
func UpdateSchema() {
	err := DbConn.AutoMigrate(
		&entities.Event{},
		&entities.User{},
		&entities.UserToken{},
		&entities.Activity{},
	)

	if err != nil {
		panic(fmt.Sprintf("Automigrate unsuccesful %v", err))
	}
}
