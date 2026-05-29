package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 允许的域名列表
		allowedOrigins := []string{
			"http://localhost", "http://localhost:5173", "http://localhost:3000",
			"http://127.0.0.1", "http://127.0.0.1:5173", "http://127.0.0.1:3000",
		}

		// 检查 origin 是否在允许列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.HasPrefix(origin, allowedOrigin) {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
