package main

import (
	"log"

	"github.com/energy-uktc/eventpool-api/controllers/api"
	"github.com/energy-uktc/eventpool-api/controllers/web"
	"github.com/energy-uktc/eventpool-api/middlewares"
	"github.com/gin-gonic/gin"
)

func configureRoutes() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/web/*.html")
	router.Static("/assets", "./assets")

	//Auth
	api.RegisterAuthRoutes(router.Group("/auth"))
	//API Groups
	v1 := router.Group("/api/v1", middlewares.AuthRequired)
	api.RegisterEventRoutes(v1.Group("/events"))
	api.RegisterUserRoutes(v1.Group("/user"))
	api.RegisterEventUserRoutes(v1.Group("/events/:eventId/users"))
	api.RegisterActivityRoutes(v1.Group("/events/:eventId/activities"))
	api.RegisterPollRoutes(v1.Group("/events/:eventId/polls"))

	//Web Groups
	webContent := router.Group("/web")
	web.RegisterAuthRoutes(webContent.Group("/auth"))
	web.RegisterEventRoutes(webContent.Group("/event/:eventId"))

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
