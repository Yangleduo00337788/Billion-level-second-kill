package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inference-engine/internal/ai"
	"inference-engine/internal/config"
	"inference-engine/internal/pkg/database"
	"inference-engine/internal/router"

	"github.com/go-redis/redis/v8"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	config.Init()
	cfg := config.Get()

	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})

	db := database.InitMySQL(&cfg.Database)

	var rdb *redis.Client
	if cfg.Redis.Host != "" {
		rdb = database.InitRedis(&cfg.Redis)
		defer rdb.Close()
	}

	aiService := ai.NewService(&cfg.AI)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
