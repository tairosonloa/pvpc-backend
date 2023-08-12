package zones

import (
	pvpc "go-pvpc/internal"
	"go-pvpc/internal/listing"
	"go-pvpc/internal/platform/storage/storagemocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_ListZonesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repositoryMock := new(storagemocks.PricesZoneRepository)
	listingService := listing.NewListingService(repositoryMock)

	r := gin.New()
	r.GET("/zones", ListZonesHandler(listingService))

	t.Run("list available zones and returns 200", func(t *testing.T) {
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

		require.Equal(t, http.StatusOK, res.StatusCode)
		snaps.MatchSnapshot(t, rec.Body.String())
	})
}
