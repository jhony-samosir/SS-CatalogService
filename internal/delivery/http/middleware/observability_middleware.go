package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CorrelationIDMiddleware ensures every request has a unique X-Correlation-Id.
// It will use the header from the gateway if present, otherwise it generates a new one.
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := ""

		// 1. Try W3C traceparent (00-traceid-spanid-flags)
		if tp := c.GetHeader("traceparent"); tp != "" {
			parts := strings.Split(tp, "-")
			if len(parts) >= 2 {
				correlationID = parts[1] // Use trace ID as correlation ID
			}
		}

		// 2. Fallback to X-Correlation-Id
		if correlationID == "" {
			correlationID = c.GetHeader("X-Correlation-Id")
		}

		// 3. Generate if still empty
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Set in context and header
		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-Id", correlationID)
		
		c.Next()
	}
}
