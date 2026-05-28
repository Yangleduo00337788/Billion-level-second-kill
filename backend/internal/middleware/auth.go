package middleware

import (
	"bytes"
	"io"
	"net"
	"strings"

	"inference-engine/internal/admin"
	"inference-engine/internal/pkg/jwt"
	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// normalizeIP converts IPv6 loopback to IPv4 format
func normalizeIP(ip string) string {
	if ip == "::1" {
		return "127.0.0.1"
	}
	// Try to parse and convert IPv6 to IPv4 if possible
	parsed := net.ParseIP(ip)
	if parsed != nil {
		if v4 := parsed.To4(); v4 != nil {
			return v4.String()
		}
	}
	return ip
}

var globalDB *gorm.DB

func SetGlobalDB(db *gorm.DB) {
	globalDB = db
}

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
		c.Set("clientIP", normalizeIP(c.ClientIP()))

		// Look up username and cache it
		if globalDB != nil && claims.UserID > 0 {
			var username string
			globalDB.Table("users").Where("id = ?", claims.UserID).Pluck("username", &username)
			if username != "" {
				c.Set("username", username)
			}
		}

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

func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

func GetRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}

func GetClientIP(c *gin.Context) string {
	ip, exists := c.Get("clientIP")
	if !exists {
		return ""
	}
	return ip.(string)
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

		c.Next()

		// Get user info
		userID := GetUserID(c)
		username := GetUsername(c)

		// Create audit log with IP
		action := c.Request.Method + " " + c.Request.URL.Path
		target := c.Request.URL.String()
		detail := string(bodyBytes)
		ip := normalizeIP(c.ClientIP())

		go handler.CreateAuditLog(userID, username, action, target, detail, ip)
	}
}

