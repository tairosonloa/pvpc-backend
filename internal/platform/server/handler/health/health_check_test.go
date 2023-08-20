package health

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HealthCheckHandlerV1_UP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	t.Run("when db is up it returns 200", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		r.GET("/v1/health", HealthCheckHandlerV1(db, 1*time.Millisecond))

		sqlMock.ExpectPing()
		req, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func Test_HealthCheckHandlerV1_DOWN(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	t.Run("when db is down it returns 503", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		r.GET("/v1/health", HealthCheckHandlerV1(db, 1*time.Millisecond))

		sqlMock.ExpectPing().WillReturnError(errors.New("mock-error"))
		req, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	})
}
