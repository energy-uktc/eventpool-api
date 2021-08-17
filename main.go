package main

import "github.com/energy-uktc/eventpool-api/database"

func main() {
	database.SetupDatabase()
	database.UpdateSchema()

	configureRoutes()
}
