package api

import (
	"time"

	"github.com/energy-uktc/grouping-api/models"
	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(g *gin.RouterGroup) {
	g.GET("", getEvent)
}

func getEvent(c *gin.Context) {
	events := make([]models.Event, 0)
	events = append(events, models.Event{Id: 1, Title: "New Event", CreatedBy: "mario.stoilov@me.com", StartDate: time.Now()})
	c.JSON(200, events)
}
