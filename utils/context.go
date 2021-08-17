package utils

import (
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/utils/constants"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) *models.UserContextInfo {
	value, ok := c.Get(constants.USER_CONTEXT)
	if !ok {
		return nil
	}
	userContext, ok := value.(models.UserContextInfo)
	if !ok {
		return nil
	}
	return &userContext
}
