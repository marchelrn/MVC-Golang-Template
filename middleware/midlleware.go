package middleware

import (
	"net/http"
	"strings"

	"mini_jira/pkg/token"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "user_id"
	ContextUsername = "username"
	ContextRole     = "role"
)

func Authenticate(tokenManager *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"message":     "token manager is not initialized",
			})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": http.StatusUnauthorized,
				"message":     "missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": http.StatusUnauthorized,
				"message":     "invalid authorization header",
			})
			c.Abort()
			return
		}

		claims, err := tokenManager.VerifyAccessToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": http.StatusUnauthorized,
				"message":     "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRole)
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"status_code": http.StatusForbidden,
				"message":     "admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
