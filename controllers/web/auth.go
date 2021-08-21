package web

import (
	"fmt"
	"net/http"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/services/user_service"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(g *gin.RouterGroup) {
	g.GET("verify", verifyCode)
	g.GET("resetPassword", resetPassword)
}

func verifyCode(c *gin.Context) {
	verificationCode := c.Query("code")
	mobileLink := c.Query("mobileLink")

	if verificationCode == "" {
		c.HTML(http.StatusOK, "web/auth/verifyUser.html", gin.H{
			"email":        "",
			"mobileLink":   fmt.Sprintf("%s?verified=%v&errorMessage=%s", mobileLink, false, "No verification provided"),
			"verified":     false,
			"errorMessage": "No verification code provided",
		})
		return
	}
	verifyUser := &models.VerifyUserRequest{
		VerificationCode:  verificationCode,
		ReturnSecureToken: true,
	}

	verifyUserResponse, err := user_service.VerifyUserCode(verifyUser)
	if err != nil {
		c.HTML(http.StatusOK, "web/auth/verifyUser.html", gin.H{
			"email":        "",
			"mobileLink":   fmt.Sprintf("%s?verified=%v&errorMessage=%s", mobileLink, false, err.Error()),
			"verified":     false,
			"errorMessage": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "web/auth/verifyUser.html", gin.H{
		"email":        verifyUserResponse.Email,
		"mobileLink":   fmt.Sprintf("%s?verified=%v", mobileLink, true),
		"verified":     true,
		"errorMessage": "",
	})
}

func resetPassword(c *gin.Context) {
	verificationCode := c.Query("code")
	mobileLink := c.Query("mobileLink")

	if verificationCode == "" {
		c.HTML(http.StatusOK, "web/auth/useMobile.html", gin.H{
			"mobileLink": fmt.Sprintf("%s?verificationCode=%v&errorMessage=%s", mobileLink, "", "No verification provided"),
		})
		return
	}

	c.HTML(http.StatusOK, "web/auth/useMobile.html", gin.H{
		"mobileLink": fmt.Sprintf("%s?verificationCode=%v", mobileLink, verificationCode),
	})
}
