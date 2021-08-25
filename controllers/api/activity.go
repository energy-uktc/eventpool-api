package api

import (
	"net/http"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/services/activity_service"
	"github.com/gin-gonic/gin"
)

func RegisterActivityRoutes(g *gin.RouterGroup) {
	g.GET("", getActivities)
	g.GET("/:activityId", getActivity)
	g.POST("", createActivity)
	g.PUT("/:activityId", updateActivity)
	g.DELETE("/:activityId", deleteActivity)
}

func getActivities(c *gin.Context) {
	eventId := c.Param("eventId")
	activities, err := activity_service.FindAll(eventId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if activities == nil {
		activities = make([]*models.Activity, 0)
	}
	c.JSON(http.StatusOK, activities)
}

func getActivity(c *gin.Context) {
	eventId := c.Param("eventId")
	activityId := c.Param("activityId")
	activity, err := activity_service.FindById(eventId, activityId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func createActivity(c *gin.Context) {
	activity := new(models.CreateUpdateActivity)
	if c.ShouldBind(activity) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}
	activity.EventID = c.Param("eventId")
	createdActivity, err := activity_service.Create(activity)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, createdActivity)
}

func updateActivity(c *gin.Context) {
	id := c.Param("activityId")
	eventId := c.Param("eventId")
	activity := new(models.CreateUpdateActivity)
	if c.ShouldBind(activity) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}
	activity.EventID = eventId
	updatedActivity, err := activity_service.Update(id, activity)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, updatedActivity)
}

func deleteActivity(c *gin.Context) {
	id := c.Param("activityId")
	err := activity_service.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
