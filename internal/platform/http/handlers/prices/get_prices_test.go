package prices

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/mocks"
	"pvpc-backend/internal/services"
	"pvpc-backend/pkg/logger"
)

func Test_GetPricesV1_Success(t *testing.T) {
	logger.SetTestLogger(os.Stderr)
	gin.SetMode(gin.TestMode)
	repositoryMock := new(mocks.PricesRepository)
	pricesService := services.NewPricesService(nil, nil, repositoryMock, nil)

	r := gin.New()
	r.GET("/v1/prices", GetPricesHandlerV1(pricesService))

	prices, err := domain.NewPrices(domain.PricesDto{
		ID:   "ABC-2023-10-02",
		Date: "2023-10-02T00:00:00+02:00",
		Zone: domain.ZoneDto{
			ID:         "ABC",
			ExternalID: "1234",
			Name:       "zone1",
		},
		Values: []domain.HourlyPriceDto{
			{
				Datetime: "2023-10-02T00:00:00+02:00",
				Value:    0.1,
			},
		},
	})
	require.NoError(t, err)

	repositoryMock.On(
		"Query",
		mock.Anything,
		(*domain.ZoneID)(nil),
		(*time.Time)(nil),
	).Return([]domain.Prices{prices}, nil)

	req, err := http.NewRequest(http.MethodGet, "/v1/prices", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	repositoryMock.AssertExpectations(t)
	require.Equal(t, http.StatusOK, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())

}

func Test_GetPricesV1_Empty(t *testing.T) {
	logger.SetTestLogger(os.Stderr)
	gin.SetMode(gin.TestMode)
	repositoryMock := new(mocks.PricesRepository)
	pricesService := services.NewPricesService(nil, nil, repositoryMock, nil)

	zoneIDRaw := "ZON"
	zoneID, err := domain.NewZoneID(zoneIDRaw)
	require.NoError(t, err)

	dateRaw := "2023-10-01"
	date, err := time.Parse("2006-01-02", dateRaw)
	require.NoError(t, err)

	r := gin.New()
	r.GET("/v1/prices", GetPricesHandlerV1(pricesService))

	repositoryMock.On(
		"Query",
		mock.Anything,
		&zoneID,
		&date,
	).Return([]domain.Prices{}, nil)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/prices?date=%s&zone_id=%s", dateRaw, zoneIDRaw), nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusNotFound, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())
}

func Test_GetPricesV1_Error(t *testing.T) {
	logger.SetTestLogger(os.Stderr)
	gin.SetMode(gin.TestMode)
	repositoryMock := new(mocks.PricesRepository)
	pricesService := services.NewPricesService(nil, nil, repositoryMock, nil)

	r := gin.New()
	r.GET("/v1/prices", GetPricesHandlerV1(pricesService))

	repositoryMock.On(
		"Query",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, errors.New("mock error"))

	req, err := http.NewRequest(http.MethodGet, "/v1/prices", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())
}
