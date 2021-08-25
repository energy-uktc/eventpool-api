package api

import (
	"net/http"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/services/poll_service"
	"github.com/energy-uktc/eventpool-api/utils"
	"github.com/gin-gonic/gin"
)

func RegisterPollRoutes(g *gin.RouterGroup) {
	g.GET("", getPolls)
	g.GET("/:pollId", getPoll)
	g.PUT("/:pollId", updateDate)
	g.POST("", createPoll)
	g.DELETE("/:pollId", deletePoll)
	g.POST("/:pollId/options/:optionId/vote", vote)
	g.DELETE("/:pollId/options/:optionId/vote", vote)
}

func getPolls(c *gin.Context) {
	eventId := c.Param("eventId")
	polls, err := poll_service.FindAll(eventId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if polls == nil {
		polls = make([]*models.PollModel, 0)
	}
	c.JSON(http.StatusOK, polls)
}

func getPoll(c *gin.Context) {
	eventId := c.Param("eventId")
	pollId := c.Param("pollId")
	poll, err := poll_service.FindById(eventId, pollId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, poll)
}

func updateDate(c *gin.Context) {

}

func createPoll(c *gin.Context) {
	poll := new(models.CreatePollModel)
	if err := c.ShouldBind(poll); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data. " + err.Error(),
		})
		return
	}
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	poll.EventID = c.Param("eventId")
	poll.CreatedBy = userContext.UserId

	createdPoll, err := poll_service.Create(poll)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, createdPoll)
}

func deletePoll(c *gin.Context) {
	id := c.Param("pollId")
	eventId := c.Param("eventId")

	err := poll_service.Delete(eventId, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func vote(c *gin.Context) {
	pollId := c.Param("pollId")
	eventId := c.Param("eventId")
	optionId := c.Param("optionId")

	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := poll_service.Vote(c.Request.Method == http.MethodPost, userContext.UserId, eventId, pollId, optionId); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
