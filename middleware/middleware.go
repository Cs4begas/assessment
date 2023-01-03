package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.GetHeader("Authorization")
		if authKey != "November 10, 2009" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
