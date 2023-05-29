package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 限流中间件 (令牌桶)
//
//	@fillInterval：添加令牌的间隔
//	@cap: 令牌桶容量
func RatelimitMiddleware(fillInterval time.Duration, cap int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(ctx *gin.Context) {
		// 如果取不到令牌就返回响应
		if bucket.TakeAvailable(1) == 0 {
			ctx.String(http.StatusOK, "rate limit...")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
