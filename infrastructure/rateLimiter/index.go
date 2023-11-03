package ratelimiter

import (
	"net/http"
	"time"

	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
	"kego.com/infrastructure/database/connection/cache"
	server_response "kego.com/infrastructure/serverResponse"
)

func RateLimiter(duration int, rate int, keySuffix string) gin.HandlerFunc {
	// Put a token into the token bucket every 1s
	// Maximum 1 request allowed per second
	return ratelimiter.GinRedisRatelimiter(cache.Client, ratelimiter.GinRatelimiterConfig{
		// config: how to generate a limit key
		LimitKey: func(c *gin.Context) string {
			return c.ClientIP() + keySuffix
		},
		// config: how to respond when limiting
		LimitedHandler: func(ctx *gin.Context) {
			server_response.Responder.Respond(ctx, http.StatusTooManyRequests, "too many requests", nil, nil)
			ctx.Abort()
		},
		// config: return ratelimiter token fill interval and bucket size
		TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
			return time.Second * time.Duration(duration), rate
		},
	})
}