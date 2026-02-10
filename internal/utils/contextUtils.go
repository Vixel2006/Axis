package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user ID not found in context. Middleware not applied or token invalid.",
		})
		return 0, fmt.Errorf("user ID not found in context")
	}

	id, ok := userID.(int)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "user ID in context is of invalid type",
		})
		return 0, fmt.Errorf("user ID in context is of invalid type")
	}

	return id, nil
}
