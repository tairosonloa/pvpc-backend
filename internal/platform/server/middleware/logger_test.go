package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	// Setting up the output recorder
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Setting up the Gin server
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(Logger([]string{}))

	// Setting up the HTTP recorder and the request
	httpRecorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test-middleware", nil)
	require.NoError(t, err)

	// Performing the request
	engine.ServeHTTP(httpRecorder, req)

	// Getting the output recorded
	require.NoError(t, w.Close())
	got, _ := io.ReadAll(r)
	os.Stderr = rescueStderr

	fmt.Println("GOT:", string(got))

	// Asserting the output contains some expected values
	assert.Contains(t, string(got), "GET")
	assert.Contains(t, string(got), "/test-middleware")
	assert.Contains(t, string(got), "404")
}
