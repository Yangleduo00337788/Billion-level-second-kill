package router

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"inference-engine/internal/admin"
	"inference-engine/internal/ai"
	"inference-engine/internal/article"
	"inference-engine/internal/chat"
	"inference-engine/internal/comment"
	"inference-engine/internal/config"
	"inference-engine/internal/middleware"
	"inference-engine/internal/notify"
	"inference-engine/internal/pkg/response"
	"inference-engine/internal/prompt"
	"inference-engine/internal/search"
	"inference-engine/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func normalizeIP(ip string) string {
	if ip == "::1" {
		return "127.0.0.1"
	}
	parsed := net.ParseIP(ip)
	if parsed != nil {
		if v4 := parsed.To4(); v4 != nil {
			return v4.String()
		}
	}
	return ip
}

func SetupRouter(db *gorm.DB, rdb *redis.Client, aiService *ai.Service) *gin.Engine {
	cfg := config.Get()
	r := gin.Default()

	middleware.SetGlobalDB(db)
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Page view recording middleware (for frontend pages)
	r.Use(middleware.PageViewRecorder(db))

	// IP 黑名单中间件
	r.Use(middleware.IPBlacklist())

	if rdb != nil {
		rateLimitCfg := middleware.NewRateLimitConfig(rdb)
		r.Use(middleware.RateLimit(rateLimitCfg))
	}

	r.GET("/api/v1/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok", "db": true, "redis": rdb != nil})
	})

	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, db)
	userHandler := user.NewHandler(userSvc)
	oauthHandler := user.NewOAuthHandler(db, &cfg.OAuth, "")
	r.GET("/api/v1/user/:id/tags", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Error(c, response.ErrBadRequest, "invalid id")
			return
		}
		var tags []admin.UserTag
		db.Where("user_id = ?", id).Find(&tags)
		response.Success(c, tags)
	})
	r.GET("/api/v1/user/list", func(c *gin.Context) {
		var users []user.User
		db.Where("status = 1").Order("article_count DESC").Limit(5).Find(&users)
		response.Success(c, users)
	})
	user.RegisterRoutes(r.Group("/api/v1"), userHandler, oauthHandler)

	articleRepo := article.NewRepository(db)
	articleSvc := article.NewService(articleRepo, userRepo, aiService, db)
	articleHandler := article.NewHandler(articleSvc)
	article.RegisterRoutes(r.Group("/api/v1"), articleHandler)

	r.GET("/api/v1/user/:id/articles", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		articles, total, _ := articleSvc.List(page, pageSize, "", 0, uint(id), "")
		response.Page(c, articles, total, page, pageSize)
	})

	r.GET("/api/v1/user/:id/favorites", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		articles, total, _ := articleSvc.GetUserFavorites(uint(id), page, pageSize)
		response.Page(c, articles, total, page, pageSize)
	})

	commentRepo := comment.NewRepository(db)
	commentSvc := comment.NewService(commentRepo, articleRepo, db)
	commentHandler := comment.NewHandler(commentSvc)
	comment.RegisterRoutes(r.Group("/api/v1"), r.Group("/api/v1/articles"), commentHandler)

	promptRepo := prompt.NewRepository(db)
	promptSvc := prompt.NewService(promptRepo, db)
	promptHandler := prompt.NewHandler(promptSvc)
	prompt.RegisterRoutes(r.Group("/api/v1"), promptHandler)

	aiHandler := ai.NewHandler(aiService)
	ai.RegisterRoutes(r.Group("/api/v1"), aiHandler)

	chatHub := chat.NewHub()
	go chatHub.Run()
	chatHandler := chat.NewHandler(chatHub)
	chat.RegisterRoutes(r.Group("/api/v1"), chatHandler)

	// Public announcements endpoint (no admin required)
	r.GET("/api/v1/announcements", func(c *gin.Context) {
		var items []admin.Announcement
		db.Where("status = ?", 1).Order("priority DESC, created_at DESC").Limit(20).Find(&items)

		// 获取当前用户ID（如果已登录）
		var userID uint
		if uid, exists := c.Get("userID"); exists {
			if id, ok := uid.(uint); ok {
				userID = id
			}
		}

		// 构建返回结果，包含已读状态
		type AnnouncementWithRead struct {
			admin.Announcement
			IsRead bool `json:"is_read"`
		}

		var results []AnnouncementWithRead
		for _, item := range items {
			isRead := false
			if userID > 0 {
				var count int64
				db.Model(&admin.AnnouncementRead{}).Where("user_id = ? AND announcement_id = ?", userID, item.ID).Count(&count)
				isRead = count > 0
			}
			results = append(results, AnnouncementWithRead{
				Announcement: item,
				IsRead:       isRead,
			})
		}

		response.Success(c, results)
	})

	// Mark announcement as read
	r.POST("/api/v1/announcements/:id/read", middleware.Auth(), func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Error(c, response.ErrBadRequest, "invalid id")
			return
		}

		userID := middleware.GetUserID(c)

		// 检查是否已读
		var count int64
		db.Model(&admin.AnnouncementRead{}).Where("user_id = ? AND announcement_id = ?", userID, id).Count(&count)
		if count == 0 {
			// 创建已读记录
			read := admin.AnnouncementRead{
				UserID:         userID,
				AnnouncementID: uint(id),
				ReadAt:         time.Now(),
			}
			db.Create(&read)
		}

		response.Success(c, nil)
	})

	// Mark all announcements as read
	r.POST("/api/v1/announcements/read-all", middleware.Auth(), func(c *gin.Context) {
		userID := middleware.GetUserID(c)

		// 获取所有公告ID
		var announcements []admin.Announcement
		db.Where("status = ?", 1).Find(&announcements)

		// 批量创建已读记录
		for _, ann := range announcements {
			var count int64
			db.Model(&admin.AnnouncementRead{}).Where("user_id = ? AND announcement_id = ?", userID, ann.ID).Count(&count)
			if count == 0 {
				read := admin.AnnouncementRead{
					UserID:         userID,
					AnnouncementID: ann.ID,
					ReadAt:         time.Now(),
				}
				db.Create(&read)
			}
		}

		response.Success(c, nil)
	})

	// Get unread announcements count
	r.GET("/api/v1/announcements/unread-count", middleware.Auth(), func(c *gin.Context) {
		userID := middleware.GetUserID(c)

		// 获取所有公告数量
		var totalAnnouncements int64
		db.Model(&admin.Announcement{}).Where("status = ?", 1).Count(&totalAnnouncements)

		// 获取已读公告数量
		var readCount int64
		db.Model(&admin.AnnouncementRead{}).Where("user_id = ?", userID).Count(&readCount)

		unreadCount := totalAnnouncements - readCount
		if unreadCount < 0 {
			unreadCount = 0
		}

		response.Success(c, gin.H{"count": unreadCount})
	})

	// Page view recording endpoint (for SPA)
	r.POST("/api/v1/page-view", func(c *gin.Context) {
		var req struct {
			Path string `json:"path" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, response.ErrBadRequest, err.Error())
			return
		}

		// Skip API paths
		if strings.HasPrefix(req.Path, "/api/") || strings.HasPrefix(req.Path, "/admin") {
			response.Success(c, nil)
			return
		}

		var userID uint
		if uid, exists := c.Get("userID"); exists {
			if id, ok := uid.(uint); ok {
				userID = id
			}
		}

		go func() {
			db.Create(&admin.PageView{
				Path:      req.Path,
				UserID:    userID,
				IP:        normalizeIP(c.ClientIP()),
				UserAgent: c.GetHeader("User-Agent"),
			})
		}()

		response.Success(c, nil)
	})

	// Public site config endpoint
	r.GET("/api/v1/site-config", func(c *gin.Context) {
		var configs []admin.SystemConfig
		db.Where("`key` IN ?", []string{"site_name", "site_description", "allow_register"}).Find(&configs)
		result := make(map[string]string)
		for _, cfg := range configs {
			result[cfg.Key] = cfg.Value
		}
		response.Success(c, result)
	})

	// User points endpoint
	r.GET("/api/v1/user/:id/points", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Error(c, response.ErrBadRequest, "invalid user id")
			return
		}
		var user struct {
			Points int `json:"points"`
		}
		db.Table("users").Select("points").Where("id = ?", id).Scan(&user)
		response.Success(c, gin.H{"points": user.Points})
	})

	// Public recommendations endpoint
	r.GET("/api/v1/recommendations", func(c *gin.Context) {
		position := c.Query("position")
		if position == "" {
			position = "homepage_top"
		}

		type RecommendResult struct {
			ID         uint   `json:"id"`
			TargetType string `json:"target_type"`
			TargetID   uint   `json:"target_id"`
			Position   string `json:"position"`
			SortOrder  int    `json:"sort_order"`
		}

		var items []RecommendResult
		db.Table("recommend_items").
			Select("id, target_type, target_id, position, sort_order").
			Where("position = ? AND status = 1", position).
			Order("sort_order ASC").
			Limit(10).
			Find(&items)

		// Fetch actual content for each recommendation
		type ContentItem struct {
			RecommendID uint   `json:"recommend_id"`
			Type        string `json:"type"`
			ID          uint   `json:"id"`
			Title       string `json:"title"`
			Summary     string `json:"summary"`
			Cover       string `json:"cover"`
			Author      string `json:"author"`
			Avatar      string `json:"avatar"`
		}

		var results []ContentItem
		for _, item := range items {
			if item.TargetType == "article" {
				var article struct {
					ID       uint   `json:"id"`
					Title    string `json:"title"`
					Summary  string `json:"summary"`
					Cover    string `json:"cover"`
					Username string `json:"username"`
					Avatar   string `json:"avatar"`
				}
				db.Table("articles").
					Select("articles.id, articles.title, articles.summary, articles.cover, users.username, users.avatar").
					Joins("LEFT JOIN users ON users.id = articles.user_id").
					Where("articles.id = ? AND articles.deleted_at IS NULL", item.TargetID).
					Scan(&article)
				if article.ID > 0 {
					results = append(results, ContentItem{
						RecommendID: item.ID,
						Type:        "article",
						ID:          article.ID,
						Title:       article.Title,
						Summary:     article.Summary,
						Cover:       article.Cover,
						Author:      article.Username,
						Avatar:      article.Avatar,
					})
				}
			} else if item.TargetType == "prompt" {
				var prompt struct {
					ID          uint   `json:"id"`
					Title       string `json:"title"`
					Description string `json:"description"`
					Username    string `json:"username"`
					Avatar      string `json:"avatar"`
				}
				db.Table("prompts").
					Select("prompts.id, prompts.title, prompts.description, users.username, users.avatar").
					Joins("LEFT JOIN users ON users.id = prompts.user_id").
					Where("prompts.id = ? AND prompts.deleted_at IS NULL", item.TargetID).
					Scan(&prompt)
				if prompt.ID > 0 {
					results = append(results, ContentItem{
						RecommendID: item.ID,
						Type:        "prompt",
						ID:          prompt.ID,
						Title:       prompt.Title,
						Summary:     prompt.Description,
						Author:      prompt.Username,
						Avatar:      prompt.Avatar,
					})
				}
			}
		}

		response.Success(c, results)
	})

	searchSvc := search.NewServiceWithDB(&cfg.ES, db)
	r.GET("/api/v1/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			response.Error(c, response.ErrBadRequest, "query is required")
			return
		}
		result, _ := searchSvc.SearchAll(query, 1, 10)
		response.Success(c, result)
	})

	notifySvc := notify.NewService(db)
	r.GET("/api/v1/notifications", middleware.Auth(), func(c *gin.Context) {
		userID := middleware.GetUserID(c)
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}
		notifications, total, _ := notifySvc.List(userID, page, pageSize)
		response.Page(c, notifications, total, page, pageSize)
	})

	r.GET("/api/v1/notifications/unread-count", middleware.Auth(), func(c *gin.Context) {
		count, _ := notifySvc.GetUnreadCount(middleware.GetUserID(c))
		response.Success(c, gin.H{"count": count})
	})

	r.PUT("/api/v1/notifications/:id/read", middleware.Auth(), func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		notifySvc.MarkAsRead(uint(id), middleware.GetUserID(c))
		response.Success(c, nil)
	})

	r.PUT("/api/v1/notifications/read-all", middleware.Auth(), func(c *gin.Context) {
		notifySvc.MarkAllAsRead(middleware.GetUserID(c))
		response.Success(c, nil)
	})

	r.POST("/api/v1/upload", middleware.Auth(), func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			response.Error(c, response.ErrBadRequest, "no file provided")
			return
		}

		// 文件大小限制 (10MB)
		if file.Size > 10*1024*1024 {
			response.Error(c, response.ErrBadRequest, "文件大小不能超过10MB")
			return
		}

		// 文件扩展名白名单
		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
			".pdf": true, ".doc": true, ".docx": true, ".md": true,
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			response.Error(c, response.ErrBadRequest, "不支持的文件类型")
			return
		}

		// 使用时间戳重命名文件，防止路径遍历
		newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst := "./uploads/" + newFilename

		if err := c.SaveUploadedFile(file, dst); err != nil {
			response.Error(c, response.ErrInternal, "文件上传失败")
			return
		}

		response.Success(c, gin.H{"url": "/uploads/" + newFilename})
	})

	adminHandler := admin.NewHandlerWithAI(db, aiService)
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Auth(), middleware.AdminOnly(), middleware.AdminAuditLog(adminHandler))
	{
		// Dashboard
		adminGroup.GET("/dashboard", adminHandler.Dashboard)
		adminGroup.GET("/chart-data", adminHandler.ChartData)
		adminGroup.GET("/page-view-stats", adminHandler.PageViewStats)

		// User Management
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.GET("/users/search", adminHandler.SearchUsers)
		adminGroup.PUT("/users/:id", adminHandler.UpdateUser)
		adminGroup.POST("/users/:id/ban", adminHandler.BanUser)
		adminGroup.POST("/users/:id/unban", adminHandler.UnbanUser)
		adminGroup.GET("/user-tags", adminHandler.ListAllUserTags)
		adminGroup.GET("/users/:id/tags", adminHandler.ListUserTags)
		adminGroup.POST("/users/:id/tags", adminHandler.AddUserTag)
		adminGroup.DELETE("/users/:id/tags/:tagId", adminHandler.DeleteUserTag)

		// Content Management
		adminGroup.GET("/articles", adminHandler.ListArticles)
		adminGroup.DELETE("/articles/:id", adminHandler.DeleteArticle)
		adminGroup.PUT("/articles/:id/moderate", adminHandler.ModerateArticle)
		adminGroup.GET("/prompts", adminHandler.ListPrompts)
		adminGroup.DELETE("/prompts/:id", adminHandler.DeletePrompt)
		adminGroup.PUT("/prompts/:id/moderate", adminHandler.ModeratePrompt)
		adminGroup.GET("/comments", adminHandler.ListComments)
		adminGroup.DELETE("/comments/:id", adminHandler.DeleteComment)

		// Content Review
		adminGroup.GET("/content-reviews", adminHandler.ListContentReviews)
		adminGroup.PUT("/content-reviews/:id", adminHandler.ReviewContent)

		// Categories & Tags
		adminGroup.POST("/categories", adminHandler.CreateCategory)
		adminGroup.DELETE("/categories/:id", adminHandler.DeleteCategory)

		// Announcements
		adminGroup.GET("/announcements", adminHandler.ListAnnouncements)
		adminGroup.POST("/announcements", adminHandler.CreateAnnouncement)
		adminGroup.PUT("/announcements/:id", adminHandler.UpdateAnnouncement)
		adminGroup.DELETE("/announcements/:id", adminHandler.DeleteAnnouncement)

		// Reports
		adminGroup.GET("/reports", adminHandler.ListReports)
		adminGroup.PUT("/reports/:id", adminHandler.HandleReport)

		// Security
		adminGroup.GET("/audit-logs", adminHandler.ListAuditLogs)
		adminGroup.GET("/login-logs", adminHandler.ListLoginLogs)
		adminGroup.GET("/ip-blacklist", adminHandler.ListIPBlacklist)
		adminGroup.POST("/ip-blacklist", adminHandler.AddIPBlacklist)
		adminGroup.DELETE("/ip-blacklist/:id", adminHandler.DeleteIPBlacklist)
		adminGroup.GET("/sensitive-words", adminHandler.ListSensitiveWords)
		adminGroup.POST("/sensitive-words", adminHandler.AddSensitiveWord)
		adminGroup.DELETE("/sensitive-words/:id", adminHandler.DeleteSensitiveWord)

		// Notifications (admin view all)
		adminGroup.GET("/notifications", adminHandler.ListAllNotifications)

		// Operations
		adminGroup.GET("/recommend-items", adminHandler.ListRecommendItems)
		adminGroup.POST("/recommend-items", adminHandler.AddRecommendItem)
		adminGroup.DELETE("/recommend-items/:id", adminHandler.DeleteRecommendItem)
		adminGroup.GET("/points-rules", adminHandler.ListPointsRules)
		adminGroup.POST("/points-rules", adminHandler.CreatePointsRule)
		adminGroup.PUT("/points-rules/:id", adminHandler.UpdatePointsRule)
		adminGroup.GET("/invite-codes", adminHandler.ListInviteCodes)
		adminGroup.POST("/invite-codes", adminHandler.CreateInviteCode)

		// System
		adminGroup.GET("/configs", adminHandler.ListConfigs)
		adminGroup.PUT("/configs/:id", adminHandler.UpdateConfig)
		adminGroup.POST("/configs", adminHandler.CreateOrUpdateConfig)
		adminGroup.GET("/ai-config", adminHandler.GetAIConfig)
		adminGroup.PUT("/ai-config", adminHandler.UpdateAIConfig)
		adminGroup.POST("/ai-config/test", adminHandler.TestAIConnection)
		adminGroup.GET("/ai-stats", adminHandler.AIUsageStats)
		adminGroup.GET("/system-logs", adminHandler.ListSystemLogs)

		// Rankings
		adminGroup.GET("/hot-articles", adminHandler.HotArticles)
		adminGroup.GET("/hot-prompts", adminHandler.HotPrompts)
		adminGroup.GET("/active-users", adminHandler.ActiveUsers)
	}

	r.POST("/api/v1/reports", middleware.Auth(), adminHandler.CreateReport)
	r.Static("/uploads", "./uploads")

	return r
}
