package database

import (
	"fmt"
	"log"
	"time"

	"inference-engine/internal/admin"
	"inference-engine/internal/article"
	"inference-engine/internal/comment"
	"inference-engine/internal/config"
	"inference-engine/internal/notify"
	"inference-engine/internal/prompt"
	"inference-engine/internal/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to mysql: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(
		// User
		&user.User{},
		&user.Follow{},

		// Article
		&article.Article{},
		&article.Category{},
		&article.Tag{},
		&article.Like{},
		&article.Favorite{},

		// Comment
		&comment.Comment{},

		// Prompt
		&prompt.Prompt{},

		// Notification
		&notify.Notification{},

		// Admin
		&admin.Announcement{},
		&admin.Report{},
		&admin.AuditLog{},
		&admin.SystemConfig{},
		&admin.AIUsageStat{},

		// Extended admin tables
		&admin.UserBan{},
		&admin.UserTag{},
		&admin.SensitiveWord{},
		&admin.ContentReview{},
		&admin.RecommendItem{},
		&admin.PointsRule{},
		&admin.InviteCode{},
		&admin.LoginLog{},
		&admin.IPBlacklist{},
		&admin.SystemLog{},
		&admin.PageView{},
	); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	// Insert default points rules if empty
	var count int64
	db.Model(&admin.PointsRule{}).Count(&count)
	if count == 0 {
		defaultRules := []admin.PointsRule{
			{Action: "register", Points: 100, Desc: "用户注册"},
			{Action: "publish_article", Points: 50, Desc: "发布文章"},
			{Action: "publish_prompt", Points: 30, Desc: "发布 Prompt"},
			{Action: "like", Points: 5, Desc: "点赞"},
			{Action: "comment", Points: 10, Desc: "评论"},
			{Action: "follow", Points: 20, Desc: "关注"},
			{Action: "daily_login", Points: 5, Desc: "每日登录"},
		}
		db.Create(&defaultRules)
	}

	// Insert default system configs if empty
	db.Model(&admin.SystemConfig{}).Count(&count)
	if count == 0 {
		defaultConfigs := []admin.SystemConfig{
			{Key: "site_name", Value: "推理引擎", Desc: "网站名称"},
			{Key: "site_description", Value: "开发者社区平台", Desc: "网站描述"},
			{Key: "allow_register", Value: "true", Desc: "是否允许注册"},
			{Key: "article_review", Value: "false", Desc: "文章是否需要审核"},
			{Key: "max_upload_size", Value: "10", Desc: "最大上传文件大小(MB)"},
			{Key: "ai_provider", Value: "openai", Desc: "AI 服务商 (openai/deepseek)"},
			{Key: "ai_api_key", Value: "", Desc: "AI API 密钥"},
			{Key: "ai_base_url", Value: "", Desc: "AI API 基础地址 (留空使用默认)"},
			{Key: "ai_model", Value: "gpt-4o", Desc: "AI 模型名称"},
			{Key: "ai_max_tokens", Value: "2048", Desc: "AI 最大 Token 数"},
			{Key: "ai_temperature", Value: "0.7", Desc: "AI Temperature"},
			{Key: "ai_rate_limit", Value: "60", Desc: "AI 每分钟请求限制"},
		}
		db.Create(&defaultConfigs)
	}

	return db
}
