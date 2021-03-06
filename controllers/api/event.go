package api

import (
	"net/http"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/services/event_service"
	"github.com/energy-uktc/eventpool-api/utils"
	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(g *gin.RouterGroup) {
	g.GET("", getEvents)
	g.GET("/:eventId", getEvent)
	g.POST("", createEvent)
	g.PUT("/:eventId", updateEvent)
	g.PATCH("/:eventId", updateEvent)
	g.DELETE("/:eventId", deleteEvent)
}

func getEvents(c *gin.Context) {
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	active := c.Query("active") == "true"
	var findEvents func(userId string) ([]*models.Event, error)
	if active {
		findEvents = event_service.FindActiveEvents
	} else {
		findEvents = event_service.FindAllForUser
	}
	events, err := findEvents(userContext.UserId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if events == nil {
		events = make([]*models.Event, 0)
	}
	c.JSON(http.StatusOK, events)
}

func getEvent(c *gin.Context) {
	id := c.Param("eventId")
	events, err := event_service.FindById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

func createEvent(c *gin.Context) {
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	event := new(models.CreateEvent)
	if c.ShouldBind(event) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}
	event.CreatedBy = userContext.UserId
	createdEvent, err := event_service.Create(event)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, createdEvent)
}

func updateEvent(c *gin.Context) {
	id := c.Param("eventId")
	event := new(models.UpdateEvent)
	if err := c.ShouldBind(event); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}
	partialUpdate := c.Request.Method == http.MethodPatch
	updatedEvent, err := event_service.Update(id, partialUpdate, event)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

func deleteEvent(c *gin.Context) {
	id := c.Param("eventId")
	err := event_service.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
