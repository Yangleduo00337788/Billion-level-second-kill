package database

import (
	"fmt"
	"log"
	"time"

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
		&user.User{},
		&user.Follow{},
		&article.Article{},
		&article.Category{},
		&article.Tag{},
		&article.Like{},
		&article.Favorite{},
		&comment.Comment{},
		&prompt.Prompt{},
		&notify.Notification{},
	); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	setTableComments(db)

	return db
}

func setTableComments(db *gorm.DB) {
	exec := db.Exec
	exec("ALTER TABLE users COMMENT '用户表'")
	exec("ALTER TABLE follows COMMENT '关注关系表'")
	exec("ALTER TABLE articles COMMENT '文章表'")
	exec("ALTER TABLE categories COMMENT '文章分类表'")
	exec("ALTER TABLE tags COMMENT '标签表'")
	exec("ALTER TABLE likes COMMENT '点赞表'")
	exec("ALTER TABLE favorites COMMENT '收藏表'")
	exec("ALTER TABLE comments COMMENT '评论表'")
	exec("ALTER TABLE prompts COMMENT 'Prompt提示词表'")
	exec("ALTER TABLE notifications COMMENT '通知表'")
}
