package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inference-engine/internal/ai"
	"inference-engine/internal/config"
	"inference-engine/internal/pkg/database"
	"inference-engine/internal/pkg/syslog"
	"inference-engine/internal/router"

	"github.com/go-redis/redis/v8"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	config.Init()
	cfg := config.Get()

	// 确保必要的目录存在
	os.MkdirAll("./uploads", 0755)
	os.MkdirAll("./logs", 0755)

	// 同时输出到文件和控制台（GoLand Run窗口可见）
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	log.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))

	db := database.InitMySQL(&cfg.Database)

	// Initialize system logger
	syslog.Init(db)
	syslog.Info("系统启动成功", "server")
	syslog.Info("数据库连接成功", "database")

	var rdb *redis.Client
	if cfg.Redis.Host != "" {
		rdb = database.InitRedis(&cfg.Redis)
		defer rdb.Close()
		syslog.Info("Redis连接成功", "redis")
	}

	aiService := ai.NewService(db)

	r := router.SetupRouter(db, rdb, aiService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		log.Printf("server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	syslog.Info("系统正在关闭", "server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
