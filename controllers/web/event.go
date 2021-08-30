package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(g *gin.RouterGroup) {
	g.GET("invite", invite)
}

func invite(c *gin.Context) {
	eventId := c.Param("eventId")
	mobileLink := c.Query("mobileLink")
	c.HTML(http.StatusOK, "web/auth/useMobile.html", gin.H{
		"mobileLink": fmt.Sprintf("%s?eventId=%v", mobileLink, eventId),
	})
}
