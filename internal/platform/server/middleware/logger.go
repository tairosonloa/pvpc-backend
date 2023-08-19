package middleware

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

// Logger is a gin.HandlerFunc that logs some information
// of the incoming request and the consequent response.
// It is intended to be used as a middleware.
func Logger(skipPaths []string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(skipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range skipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Prepare (pre-request)
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Results (post-request)
		if _, ok := skip[path]; !ok {
			latency := time.Since(start).Truncate(time.Millisecond)
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if c.Request.URL.RawQuery != "" {
				path = path + "?" + c.Request.URL.RawQuery
			}

			if errorMessage != "" {
				log.Error("Received request",
					"statusCode", statusCode,
					"latency", latency,
					"clientIP", clientIP,
					"method", method,
					"path", path,
					"error", errorMessage,
				)
			} else {
				log.Info("Received request",
					"statusCode", statusCode,
					"latency", latency,
					"clientIP", clientIP,
					"method", method,
					"path", path,
				)
			}
		}
	}
}
