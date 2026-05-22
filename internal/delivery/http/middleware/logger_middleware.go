package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware logs HTTP requests in structured JSON format.
func RequestLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start)

		// Extract path, fall back to URL Path if FullPath is not defined
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Log request completion. Use InfoContext to propagate correlation ID from standard context.
		logger.InfoContext(c.Request.Context(), "http_request",
			"method",      c.Request.Method,
			"path",        path,
			"status",      c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"service",     "ss-catalog-service",
		)
	}
}
