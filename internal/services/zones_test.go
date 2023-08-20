package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/errors"
	"pvpc-backend/internal/platform/storage/storagemocks"
)

func Test_ZonesService_ListZones(t *testing.T) {

	t.Run("fails with a repository error", func(t *testing.T) {
		repositoryMock := new(storagemocks.PricesZonesRepository)
		repositoryMock.On("GetAll", mock.Anything).Return(nil, errors.NewDomainError(errors.PersistenceError, "mock-error"))

		listingService := NewZonesService(repositoryMock)
		res, err := listingService.ListZones(context.Background())
		require.Error(t, err)
		assert.Nil(t, res)

		repositoryMock.AssertExpectations(t)
	})

	t.Run("succeeds and returns PricesZone's", func(t *testing.T) {
		zone, err := domain.NewPricesZone(domain.PricesZoneDto{
			ID:         "ZON",
			ExternalID: "123",
			Name:       "Zone 1",
		})
		require.NoError(t, err)

		repositoryMock := new(storagemocks.PricesZonesRepository)
		repositoryMock.On("GetAll", mock.Anything).Return([]domain.PricesZone{zone}, nil)

		listingService := NewZonesService(repositoryMock)
		res, err := listingService.ListZones(context.Background())
		require.NoError(t, err)
		assert.Equal(t, zone, res[0])

		repositoryMock.AssertExpectations(t)
	})
}
