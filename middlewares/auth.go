package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/energy-uktc/grouping-api/models"
	"github.com/energy-uktc/grouping-api/services/jwt_service"
	"github.com/energy-uktc/grouping-api/utils/constants"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(tokenString, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt_service.VerifyToken(tokenString, false)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenClaims := jwt_service.GetClaims(token)
	c.Set(constants.USER_CONTEXT, models.UserContextInfo{UserId: tokenClaims.CustomerInfo.Id, Scopes: tokenClaims.Scopes})
}
