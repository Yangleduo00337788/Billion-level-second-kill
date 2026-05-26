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
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	if rdb != nil {
		rateLimitCfg := middleware.NewRateLimitConfig(rdb)
		r.Use(middleware.RateLimit(rateLimitCfg))
	}

	r.GET("/api/v1/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
			"db":     true,
			"redis":  rdb != nil,
		})
	})

	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, db)
	userHandler := user.NewHandler(userSvc)
	user.RegisterRoutes(r.Group("/api/v1"), userHandler)

	articleRepo := article.NewRepository(db)
	articleSvc := article.NewService(articleRepo, userRepo, aiService, db)
	articleHandler := article.NewHandler(articleSvc)
	article.RegisterRoutes(r.Group("/api/v1"), articleHandler)

	r.GET("/api/v1/user/:id/articles", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Error(c, response.ErrBadRequest, "invalid user id")
			return
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 50 {
			pageSize = 10
		}
		articles, total, err := articleSvc.List(page, pageSize, "", 0, uint(id))
		if err != nil {
			response.Error(c, response.ErrInternal, err.Error())
			return
		}
		response.Page(c, articles, total, page, pageSize)
	})

	r.GET("/api/v1/user/:id/favorites", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Error(c, response.ErrBadRequest, "invalid user id")
			return
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 50 {
			pageSize = 10
		}
		articles, total, err := articleSvc.GetUserFavorites(uint(id), page, pageSize)
		if err != nil {
			response.Error(c, response.ErrInternal, err.Error())
			return
		}
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

	searchSvc := search.NewService(&config.Get().ES)
	r.GET("/api/v1/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			response.Error(c, response.ErrBadRequest, "query is required")
			return
		}

		p := 1
		ps := 10

		result, err := searchSvc.SearchAll(query, p, ps)
		if err != nil {
			response.Error(c, response.ErrInternal, err.Error())
			return
		}

		response.Success(c, result)
	})

	notifySvc := notify.NewService(db)

	r.GET("/api/v1/notifications", middleware.Auth(), func(c *gin.Context) {
		userID := middleware.GetUserID(c)
		notifications, total, err := notifySvc.List(userID, 1, 20)
		if err != nil {
			response.Error(c, response.ErrInternal, err.Error())
			return
		}
		response.Page(c, notifications, total, 1, 20)
	})

	r.GET("/api/v1/user/list", func(c *gin.Context) {
		var users []user.User
		db.Where("status = 1").Order("article_count DESC").Limit(5).Find(&users)
		response.Success(c, users)
	})

	r.POST("/api/v1/upload", middleware.Auth(), func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			response.Error(c, response.ErrBadRequest, "no file provided")
			return
		}

		dst := "./uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			response.Error(c, response.ErrInternal, err.Error())
			return
		}

		response.Success(c, gin.H{
			"url": "/uploads/" + file.Filename,
		})
	})

	adminHandler := admin.NewHandler(db)
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.Auth())
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/dashboard", adminHandler.Dashboard)
		admin.GET("/users", adminHandler.ListUsers)
		admin.PUT("/users/:id", adminHandler.UpdateUser)
		admin.GET("/articles", adminHandler.ListArticles)
		admin.DELETE("/articles/:id", adminHandler.DeleteArticle)
		admin.GET("/prompts", adminHandler.ListPrompts)
		admin.DELETE("/prompts/:id", adminHandler.DeletePrompt)
	}

	r.Static("/uploads", "./uploads")

	return r
}
