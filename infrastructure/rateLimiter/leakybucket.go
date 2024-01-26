package ratelimiter

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var rl = ratelimit.New(200) // 200 requests allwoed per second

func LeakyBucket() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rl.Take()
	}
}