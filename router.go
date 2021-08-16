package main

import (
	"log"

	"github.com/energy-uktc/grouping-api/controllers/api"
	"github.com/energy-uktc/grouping-api/controllers/web"
	"github.com/energy-uktc/grouping-api/middlewares"
	"github.com/gin-gonic/gin"
)

func configureRoutes() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/web/**/*")
	//Auth
	api.RegisterAuthRoutes(router.Group("/auth"))
	//API Groups
	v1 := router.Group("/api/v1", middlewares.AuthRequired)
	api.RegisterEventRoutes(v1.Group("/events"))

	webContent := router.Group("/web")
	web.RegisterAuthRoutes(webContent.Group("/auth"))
	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
