package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/errors"
	"pvpc-backend/internal/mocks"
)

func Test_ZonesService_ListZones(t *testing.T) {

	t.Run("fails with a repository error", func(t *testing.T) {
		repositoryMock := new(mocks.ZonesRepository)
		repositoryMock.On("GetAll", mock.Anything).Return(nil, errors.NewDomainError(errors.PersistenceError, "mock-error"))

		listingService := NewZonesService(repositoryMock)
		res, err := listingService.ListZones(context.Background())
		require.Error(t, err)
		assert.Nil(t, res)

		repositoryMock.AssertExpectations(t)
	})

	t.Run("succeeds and returns Zone's", func(t *testing.T) {
		zone, err := domain.NewZone(domain.ZoneDto{
			ID:         "ZON",
			ExternalID: "123",
			Name:       "Zone 1",
		})
		require.NoError(t, err)

		repositoryMock := new(mocks.ZonesRepository)
		repositoryMock.On("GetAll", mock.Anything).Return([]domain.Zone{zone}, nil)

		listingService := NewZonesService(repositoryMock)
		res, err := listingService.ListZones(context.Background())
		require.NoError(t, err)
		assert.Equal(t, zone, res[0])

		repositoryMock.AssertExpectations(t)
	})
}
