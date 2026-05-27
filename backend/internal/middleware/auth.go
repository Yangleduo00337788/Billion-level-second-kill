package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"inference-engine/internal/admin"
	"inference-engine/internal/pkg/jwt"
	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithStatus(c, 401, response.ErrUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithStatus(c, 401, response.ErrUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.ErrorWithStatus(c, 401, response.ErrUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(uint)
}

func GetRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c)
		if role != "admin" {
			response.ErrorWithStatus(c, 403, response.ErrForbidden, "admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminAuditLog middleware for admin write operations
func AdminAuditLog(handler *admin.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip GET requests
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Get user info
		userID := GetUserID(c)
		username := ""
		if u, exists := c.Get("username"); exists {
			username = u.(string)
		}

		// Create audit log
		action := c.Request.Method + " " + c.Request.URL.Path
		target := c.Request.URL.String()
		detail := string(bodyBytes)
		if duration > 0 {
			detail += " | duration: " + time.Since(start).String()
		}

		go handler.CreateAuditLog(userID, username, action, target, detail, c.ClientIP())
	}
}
