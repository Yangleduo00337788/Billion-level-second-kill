package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PageView struct {
	ID        uint      `gorm:"primarykey"`
	Path      string    `gorm:"size:200;not null"`
	UserID    uint      `gorm:"index"`
	IP        string    `gorm:"size:50"`
	UserAgent string    `gorm:"size:500"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (PageView) TableName() string {
	return "page_views"
}

func PageViewRecorder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		path := c.Request.URL.Path
		// Only record page visits, skip API and static assets
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/uploads/") ||
			path == "/favicon.ico" {
			return
		}

		// Skip if it's an API-only request
		accept := c.GetHeader("Accept")
		if strings.Contains(accept, "application/json") && !strings.Contains(accept, "text/html") {
			return
		}

		var userID uint
		if uid, exists := c.Get("userID"); exists {
			if id, ok := uid.(uint); ok {
				userID = id
			}
		}

		go func() {
			db.Create(&PageView{
				Path:      path,
				UserID:    userID,
				IP:        normalizeIP(c.ClientIP()),
				UserAgent: c.GetHeader("User-Agent"),
			})
		}()
	}
}