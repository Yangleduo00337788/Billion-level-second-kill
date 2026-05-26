package middleware

import (
	"context"
	"net/http"
	"time"

	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimitConfig struct {
	RDB     *redis.Client
	Rate    int
	Burst   int
	Expires time.Duration
}

func NewRateLimitConfig(rdb *redis.Client) *RateLimitConfig {
	return &RateLimitConfig{
		RDB:     rdb,
		Rate:    10,
		Burst:   20,
		Expires: time.Second,
	}
}

func RateLimit(cfg *RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip
		ctx := context.Background()

		pipe := cfg.RDB.Pipeline()
		countCmd := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, cfg.Expires)
		_, err := pipe.Exec(ctx)
		if err != nil {
			c.Next()
			return
		}

		if countCmd.Val() > int64(cfg.Burst) {
			response.ErrorWithStatus(c, http.StatusTooManyRequests, response.ErrBadRequest, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
