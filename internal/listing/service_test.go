package listing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pvpc "go-pvpc/internal"
	"go-pvpc/internal/errors"
	"go-pvpc/internal/platform/storage/storagemocks"
)

func Test_ListingService(t *testing.T) {

	t.Run("fails with a repository error", func(t *testing.T) {
		repositoryMock := new(storagemocks.PricesZonesRepository)
		repositoryMock.On("GetAll", mock.Anything).Return(nil, errors.NewDomainError(errors.PersistenceError, "mock-error"))

		listingService := NewListingService(repositoryMock)
		res, err := listingService.ListPricesZones(context.Background())
		require.Error(t, err)
		assert.Nil(t, res)

		repositoryMock.AssertExpectations(t)
	})

	t.Run("succeeds and returns PricesZone's", func(t *testing.T) {
		zone, err := pvpc.NewPricesZone(pvpc.PricesZoneDto{
			ID:         "ZON",
			ExternalID: "123",
			Name:       "Zone 1",
		})
		require.NoError(t, err)

		repositoryMock := new(storagemocks.PricesZonesRepository)
		repositoryMock.On("GetAll", mock.Anything).Return([]pvpc.PricesZone{zone}, nil)

		listingService := NewListingService(repositoryMock)
		res, err := listingService.ListPricesZones(context.Background())
		require.NoError(t, err)
		assert.Equal(t, zone, res[0])

		repositoryMock.AssertExpectations(t)
	})
}
