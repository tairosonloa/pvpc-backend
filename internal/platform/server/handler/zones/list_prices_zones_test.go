package zones

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pvpc "pvpc-backend/internal"
	"pvpc-backend/internal/listing"
	"pvpc-backend/internal/platform/storage/storagemocks"
)

func Test_ListZonesHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repositoryMock := new(storagemocks.PricesZonesRepository)
	listingService := listing.NewListingService(repositoryMock)

	r := gin.New()
	r.GET("/zones", ListZonesHandler(listingService))

	zone1, err := pvpc.NewPricesZone(pvpc.PricesZoneDto{
		ID:         "ABC",
		ExternalID: "1234",
		Name:       "zone1",
	})
	require.NoError(t, err)

	zone2, err := pvpc.NewPricesZone(pvpc.PricesZoneDto{
		ID:         "DEF",
		ExternalID: "5678",
		Name:       "zone2",
	})
	require.NoError(t, err)

	repositoryMock.On(
		"GetAll",
		mock.Anything,
	).Return([]pvpc.PricesZone{zone1, zone2}, nil)

	req, err := http.NewRequest(http.MethodGet, "/zones", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	repositoryMock.AssertExpectations(t)
	require.Equal(t, http.StatusOK, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())

}

func Test_ListZonesHandler_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repositoryMock := new(storagemocks.PricesZonesRepository)
	listingService := listing.NewListingService(repositoryMock)

	r := gin.New()
	r.GET("/zones", ListZonesHandler(listingService))

	repositoryMock.On(
		"GetAll",
		mock.Anything,
	).Return([]pvpc.PricesZone{}, nil)

	req, err := http.NewRequest(http.MethodGet, "/zones", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())
}

func Test_ListZonesHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repositoryMock := new(storagemocks.PricesZonesRepository)
	listingService := listing.NewListingService(repositoryMock)

	r := gin.New()
	r.GET("/zones", ListZonesHandler(listingService))

	repositoryMock.On(
		"GetAll",
		mock.Anything,
	).Return(nil, errors.New("mock error"))

	req, err := http.NewRequest(http.MethodGet, "/zones", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	snaps.MatchSnapshot(t, rec.Body.String())

}
