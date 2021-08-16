package web

import (
	"net/http"

	"github.com/energy-uktc/grouping-api/models"
	"github.com/energy-uktc/grouping-api/services/user_service"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(g *gin.RouterGroup) {
	g.GET("verify", verifyCode)
}

func verifyCode(c *gin.Context) {
	verificationCode := c.Query("code")
	if verificationCode == "" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"errorMessage": "No verification provided",
		})
		return
	}
	mobileLink := c.Query("mobileLink")
	verifyUser := &models.VerifyUserRequest{
		VerificationCode:  verificationCode,
		ReturnSecureToken: true,
	}

	verifyUserResponse, err := user_service.VerifyUserCode(verifyUser)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "web/auth/verifyUser.tmpl", gin.H{
		"email":      verifyUserResponse.Email,
		"mobileLink": mobileLink,
	})
}
