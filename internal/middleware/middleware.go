package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()
		method := ctx.Request.Method
		clientIP := ctx.ClientIP()

		attrs := []any{
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Duration("latency", duration),
		}

		if raw != "" {
			attrs = append(attrs, slog.String("query", raw))
		}

		if status >= 500 {
			slog.Error("Request failed", attrs...)
		} else if status >= 400 {
			slog.Warn("Bad request", attrs...)
		} else {
			slog.Info("Request success", attrs...)
		}
	}
}
