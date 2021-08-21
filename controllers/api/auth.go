package api

import (
	"log"
	"net/http"

	"github.com/energy-uktc/eventpool-api/middlewares"
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/services/mail_service"
	"github.com/energy-uktc/eventpool-api/services/user_service"
	"github.com/energy-uktc/eventpool-api/utils"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(g *gin.RouterGroup) {
	g.POST("register", registerUser)
	g.POST("token", generateToken)
	g.DELETE("token", middlewares.AuthRequired, revokeRefreshToken)
	g.POST("verify", verifyCode)
	g.POST("resendVerificationCode", resendVerificationCode)
	g.POST("refreshToken", refreshToken)
	g.POST("password", changePassword)
}

func registerUser(c *gin.Context) {
	user := new(models.CreateUserRequest)
	if c.Bind(user) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	mobileAppUrl := c.Query("mobileAppUrl")
	createdUser, err := user_service.RegisterUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	verificationCode, err := user_service.GenerateVerificationCode(createdUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	go mail_service.VerificationCodeRequest(createdUser.UserName, createdUser.Email, verificationCode, mobileAppUrl)
	c.Status(http.StatusCreated)
}

func resendVerificationCode(c *gin.Context) {
	authUserRequest := new(models.AuthUserRequest)
	if c.Bind(authUserRequest) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	mobileAppUrl := c.Query("mobileAppUrl")

	user, err := user_service.FindUnverifiedUser(authUserRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	verificationCode, err := user_service.GenerateVerificationCode(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	go mail_service.VerificationCodeRequest(user.UserName, user.Email, verificationCode, mobileAppUrl)
	c.Status(http.StatusOK)
}

func verifyCode(c *gin.Context) {
	verifyUser := new(models.VerifyUserRequest)
	if c.Bind(verifyUser) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	verifyUserResponse, err := user_service.VerifyUserCode(verifyUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if verifyUserResponse != nil {
		c.JSON(http.StatusOK, verifyUserResponse.Token)
		return
	}
	c.Status(http.StatusOK)
}

func generateToken(c *gin.Context) {
	user := new(models.AuthUserRequest)
	if err := c.Bind(user); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := user_service.GenerateToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, token)
}

func refreshToken(c *gin.Context) {
	refreshTokenRequest := new(models.RefreshTokenRequest)
	if c.Bind(refreshTokenRequest) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := user_service.RefreshToken(refreshTokenRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, token)
}

func revokeRefreshToken(c *gin.Context) {
	userContext := utils.GetUserFromContext(c)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err := user_service.RevokeRefreshToken(userContext)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func changePassword(c *gin.Context) {
	changePasswordRequest := new(models.ChangePasswordRequest)
	if c.Bind(changePasswordRequest) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	switch changePasswordRequest.Action {
	case "change":
		if err := user_service.ChangePassword(changePasswordRequest.Email, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	case "reset":
		if err := user_service.ResetPassword(changePasswordRequest.VerificationCode, changePasswordRequest.NewPassword); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	case "sendResetEmail":
		mobileAppUrl := c.Query("mobileAppUrl")
		user, err := user_service.FindUserByEmail(changePasswordRequest.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if code, err := user_service.GenerateResetPasswordCode(user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			go mail_service.ResetPasswordRequest(user.UserName, user.Email, code, mobileAppUrl)
		}
	}

	c.Status(http.StatusOK)
}
