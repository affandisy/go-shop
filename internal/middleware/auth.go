package middleware

import (
	"strings"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/affandisy/goshop/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "User role not found")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", domain.ErrUnauthorized
	}

	return userID.(string), nil
}

func GetUserRole(c *gin.Context) (string, error) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", domain.ErrUnauthorized
	}

	return role.(string), nil
}
