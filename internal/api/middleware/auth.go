package middleware

import (
	"av-merch-shop/internal/api/common"
	"av-merch-shop/internal/entities"
	"av-merch-shop/pkg/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "authorization header required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, "Bearer ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "invalid token"})
			c.Abort()
			return
		}

		c.Set(common.ContextUserID, claims.UserID)
		c.Set(common.ContextRole, claims.Role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(common.ContextRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "unauthorized"})
			c.Abort()
			return
		}

		if role != entities.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"errors": "admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
