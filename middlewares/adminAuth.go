package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context (hasil AuthJWT)
		roleInterface, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - role not found"})
			return
		}

		role, ok := roleInterface.(string)
		if !ok || role != "Admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Admin access only"})
			return
		}

		c.Next()
	}
}
