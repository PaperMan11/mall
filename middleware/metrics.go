package middleware

import (
	"mall/pkg/utils/metrics"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RecordMetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		cost := time.Since(start).Seconds()
		metrics.RecordMetrics(
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Writer.Status() == http.StatusOK,
			ctx.Writer.Status(),
			cost,
		)
	}
}
