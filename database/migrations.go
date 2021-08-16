package database

import (
	"fmt"

	"github.com/energy-uktc/grouping-api/entities"
)

//UpdateSchema performs database migrations for used entities
func UpdateSchema() {
	err := DbConn.AutoMigrate(
		&entities.Event{},
		&entities.User{},
		&entities.UserToken{},
	)

	if err != nil {
		panic(fmt.Sprintf("Automigrate unsuccesful %v", err))
	}
}
