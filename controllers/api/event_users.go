package api

import (
	"net/http"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/services/event_service"
	"github.com/energy-uktc/eventpool-api/services/mail_service"
	"github.com/energy-uktc/eventpool-api/services/user_service"
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

	user, err := user_service.GetUser(userContext.UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	eventId := c.Param("eventId")
	invitation := c.Query("invitation") == "true"

	if invitation {
		mobileAppUrl := c.Query("mobileAppUrl")
		invitationRequest := new(models.InvitationRequest)
		if err := c.ShouldBind(invitationRequest); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Data",
			})
			return
		}

		event, err := event_service.FindById(invitationRequest.EventId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		for _, email := range invitationRequest.Emails {
			invitedUser, _ := user_service.FindUserByEmail(email)
			invitee := email
			if invitedUser != nil {
				invitee = invitedUser.UserName
			}
			go mail_service.SendInvitation(invitee, user.UserName, event.Id, event.Title, email, mobileAppUrl)
		}
	} else {
		err := event_service.AssignUser(eventId, userContext.UserId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
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
