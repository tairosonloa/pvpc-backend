package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pvpc-backend/pkg/logger"
)

func TestMiddleware(t *testing.T) {
	// Setting up the output recorder
	r, w, _ := os.Pipe()

	// Setting up the Gin server
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(Logger([]string{}))
	logger.SetTestLogger(w)

	// Setting up the HTTP recorder and the request
	httpRecorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test-middleware", nil)
	require.NoError(t, err)

	// Performing the request
	engine.ServeHTTP(httpRecorder, req)

	// Getting the output recorded
	require.NoError(t, w.Close())
	got, _ := io.ReadAll(r)

	// Asserting the output contains some expected values
	assert.Contains(t, string(got), "GET")
	assert.Contains(t, string(got), "/test-middleware")
	assert.Contains(t, string(got), "404")
}
