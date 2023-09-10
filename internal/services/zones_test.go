package services

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/domain/errors"
	"pvpc-backend/internal/mocks"
	"pvpc-backend/pkg/logger"
)

func Test_ZonesService_ListZones(t *testing.T) {
	logger.SetTestLogger(os.Stderr)

	t.Run("fails with a repository error", func(t *testing.T) {
		repositoryMock := new(mocks.ZonesRepository)
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		repositoryMock.On("GetAll", mock.Anything).Return(nil, mockError)

		zonesService := NewZonesService(repositoryMock)
		res, err := zonesService.ListZones(context.Background())
		require.Error(t, err)
		require.Equal(t, mockError, err)
		require.Nil(t, res)

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

		zonesService := NewZonesService(repositoryMock)
		res, err := zonesService.ListZones(context.Background())
		require.NoError(t, err)
		require.Equal(t, zone, res[0])

		repositoryMock.AssertExpectations(t)
	})
}
