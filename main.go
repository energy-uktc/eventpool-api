package main

import "github.com/energy-uktc/grouping-api/database"

func main() {
	database.SetupDatabase()
	database.UpdateSchema()

	configureRoutes()
}
