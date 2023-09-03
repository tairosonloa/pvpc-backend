package zones

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/mocks"
	"pvpc-backend/internal/services"
	"pvpc-backend/pkg/logger"
)

func Test_ListZonesHandlerV1_Success(t *testing.T) {
	logger.SetTestLogger(os.Stderr)
	gin.SetMode(gin.TestMode)
	repositoryMock := new(mocks.ZonesRepository)
	listingService := services.NewZonesService(repositoryMock)

	r := gin.New()
	r.GET("/v1/zones", ListZonesHandlerV1(listingService))

	zone1, err := domain.NewZone(domain.ZoneDto{
		ID:         "ABC",
		ExternalID: "1234",
		Name:       "zone1",
	})
	require.NoError(t, err)

	zone2, err := domain.NewZone(domain.ZoneDto{
		ID:         "DEF",
		ExternalID: "5678",
		Name:       "zone2",
	})
	require.NoError(t, err)

	repositoryMock.On(
		"GetAll",
		mock.Anything,
	).Return([]domain.Zone{zone1, zone2}, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/zones", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	repositoryMock.AssertExpectations(t)
	require.Equal(t, http.StatusOK, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())

}

func Test_ListZonesHandlerV1_Empty(t *testing.T) {
	logger.SetTestLogger(os.Stderr)
	gin.SetMode(gin.TestMode)
	repositoryMock := new(mocks.ZonesRepository)
	listingService := services.NewZonesService(repositoryMock)

	r := gin.New()
	r.GET("/v1/zones", ListZonesHandlerV1(listingService))

	repositoryMock.On(
		"GetAll",
		mock.Anything,
	).Return([]domain.Zone{}, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/zones", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())
}

func Test_ListZonesHandlerV1_Error(t *testing.T) {
	logger.SetTestLogger(os.Stderr)
	gin.SetMode(gin.TestMode)
	repositoryMock := new(mocks.ZonesRepository)
	listingService := services.NewZonesService(repositoryMock)

	r := gin.New()
	r.GET("/v1/zones", ListZonesHandlerV1(listingService))

	repositoryMock.On(
		"GetAll",
		mock.Anything,
	).Return(nil, errors.New("mock error"))

	req, err := http.NewRequest(http.MethodGet, "/v1/zones", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())

}
