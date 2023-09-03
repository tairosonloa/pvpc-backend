package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pvpc-backend/internal/platform/http/responses"
	"pvpc-backend/pkg/logger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

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
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Results (post-request)
		if _, ok := skip[path]; !ok {
			latency := time.Since(start).Truncate(time.Millisecond)
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()

			if c.Request.URL.RawQuery != "" {
				path = path + "?" + c.Request.URL.RawQuery
			}

			if statusCode >= http.StatusBadRequest {
				response := responses.APIErrorResponse{}
				json.Unmarshal(blw.body.Bytes(), &response)

				logger.Error("Errored request",
					"statusCode", statusCode,
					"latency", latency,
					"clientIP", clientIP,
					"method", method,
					"path", path,
					"errCode", response.ErrorCode,
					"errMsg", response.Message,
				)
			} else {
				logger.Info("Handled request",
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
