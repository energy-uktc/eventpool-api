package api

import (
	"net/http"

	"github.com/energy-uktc/eventpool-api/services/event_service"
	"github.com/energy-uktc/eventpool-api/utils"
	"github.com/gin-gonic/gin"
)

func RegisterEventUserRoutes(g *gin.RouterGroup) {
	g.POST("", joinEvent)
	g.DELETE("", leaveEvent)
}

func joinEvent(c *gin.Context) {
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	eventId := c.Param("eventId")
	err := event_service.AssignUser(eventId, userContext.UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func leaveEvent(c *gin.Context) {
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	eventId := c.Param("eventId")
	err := event_service.RemoveUser(eventId, userContext.UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
