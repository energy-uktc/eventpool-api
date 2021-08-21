package api

import (
	"net/http"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/services/user_service"
	"github.com/energy-uktc/eventpool-api/utils"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(g *gin.RouterGroup) {
	g.GET("", getUser)
	g.PUT("", updateUser)
}

func getUser(c *gin.Context) {
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user, err := user_service.GetUser(userContext.UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	updatedUser := new(models.UpdateUserModel)
	if c.Bind(updatedUser) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user data",
		})
		return
	}
	user, err := user_service.UpdateUser(userContext.UserId, updatedUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
