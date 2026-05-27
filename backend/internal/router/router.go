package router

import (
	"strconv"

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

func SetupRouter(db *gorm.DB, rdb *redis.Client, aiService *ai.Service) *gin.Engine {
	cfg := config.Get()
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

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
		articles, total, _ := articleSvc.List(page, pageSize, "", 0, uint(id))
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
	promptSvc := prompt.NewService(promptRepo)
	promptHandler := prompt.NewHandler(promptSvc)
	prompt.RegisterRoutes(r.Group("/api/v1"), promptHandler)

	aiHandler := ai.NewHandler(aiService)
	ai.RegisterRoutes(r.Group("/api/v1"), aiHandler)

	chatHub := chat.NewHub()
	go chatHub.Run()
	chatHandler := chat.NewHandler(chatHub)
	chat.RegisterRoutes(r.Group("/api/v1"), chatHandler)

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
		notifications, total, _ := notifySvc.List(userID, 1, 20)
		response.Page(c, notifications, total, 1, 20)
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
		dst := "./uploads/" + file.Filename
		c.SaveUploadedFile(file, dst)
		response.Success(c, gin.H{"url": "/uploads/" + file.Filename})
	})

	adminHandler := admin.NewHandler(db)
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Auth(), middleware.AdminOnly(), middleware.AdminAuditLog(adminHandler))
	{
		// Dashboard
		adminGroup.GET("/dashboard", adminHandler.Dashboard)
		adminGroup.GET("/chart-data", adminHandler.ChartData)
		adminGroup.GET("/page-view-stats", adminHandler.PageViewStats)

		// User Management
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.PUT("/users/:id", adminHandler.UpdateUser)
		adminGroup.POST("/users/:id/ban", adminHandler.BanUser)
		adminGroup.POST("/users/:id/unban", adminHandler.UnbanUser)
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

		// Operations
		adminGroup.GET("/recommend-items", adminHandler.ListRecommendItems)
		adminGroup.POST("/recommend-items", adminHandler.AddRecommendItem)
		adminGroup.DELETE("/recommend-items/:id", adminHandler.DeleteRecommendItem)
		adminGroup.GET("/points-rules", adminHandler.ListPointsRules)
		adminGroup.PUT("/points-rules/:id", adminHandler.UpdatePointsRule)
		adminGroup.GET("/invite-codes", adminHandler.ListInviteCodes)
		adminGroup.POST("/invite-codes", adminHandler.CreateInviteCode)

		// System
		adminGroup.GET("/configs", adminHandler.ListConfigs)
		adminGroup.PUT("/configs/:id", adminHandler.UpdateConfig)
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
