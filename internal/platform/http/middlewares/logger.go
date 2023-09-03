package middlewares

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}
		method := c.Request.Method
		clientIP := c.ClientIP()
		requestID := uuid.New().String()
		c.Set(logger.ContextKeyRequestID, requestID)

		logger.LogAttrs(c, slog.LevelInfo, "Received new request",
			slog.String("clientIP", clientIP),
			slog.String("method", method),
			slog.String("path", path),
		)

		// Process request
		c.Next()

		// Results (post-request)
		if _, ok := skip[path]; !ok {
			latency := time.Since(start).Truncate(time.Millisecond).String()
			statusCode := c.Writer.Status()

			if statusCode >= http.StatusBadRequest {
				response := responses.APIErrorResponse{}
				json.Unmarshal(blw.body.Bytes(), &response)

				logger.LogAttrs(c, slog.LevelError, "Errored request",
					slog.String("clientIP", clientIP),
					slog.String("method", method),
					slog.String("path", path),
					slog.Int("status", statusCode),
					slog.String("latency", latency),
					slog.String("errCode", response.ErrorCode),
					slog.String("errMsg", response.Message),
				)
			} else {
				logger.LogAttrs(c, slog.LevelInfo, "Handled request",
					slog.String("client_ip", clientIP),
					slog.String("method", method),
					slog.String("path", path),
					slog.Int("status", statusCode),
					slog.String("latency", latency),
				)
			}
		}
	}
}
