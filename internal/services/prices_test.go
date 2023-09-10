package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/domain/errors"
	"pvpc-backend/internal/mocks"
	"pvpc-backend/pkg/logger"
)

func restoreNow(nowFunc func() time.Time) {
	now = nowFunc
}

func Test_PricesService_FetchAndStorePricesFromREE(t *testing.T) {
	logger.SetTestLogger(os.Stderr)
	testZoneDto := domain.ZoneDto{ID: "ZON", ExternalID: "123", Name: "Zone 1"}
	testValuesDto := domain.HourlyPriceDto{Datetime: "2020-01-01T00:00:00Z", Value: 0.123}
	testPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: "2020-01-01T00:00:00Z", Values: []domain.HourlyPriceDto{testValuesDto}})
	require.NoError(t, err)
	testPricesFetchIdStr := "ZON-2023-01-01"
	testPricesFetchId, err := domain.NewPricesID(testPricesFetchIdStr)
	require.NoError(t, err)
	testPricesFetch, err := domain.NewPrices(domain.PricesDto{ID: testPricesFetchIdStr, Zone: testZoneDto, Date: "2023-01-01T00:00:00Z", Values: []domain.HourlyPriceDto{testValuesDto}})
	require.NoError(t, err)
	testZone, err := domain.NewZone(testZoneDto)
	require.NoError(t, err)
	todayTestDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	tomorrowTestDate := todayTestDate.AddDate(0, 0, 1)
	now = func() time.Time { return todayTestDate }
	defer restoreNow(time.Now)

	t.Run("fails with a repository error getting prices", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return(nil, mockError)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.Error(t, err)
		require.Equal(t, mockError, err)
		require.Nil(t, res)

		pricesRepositoryMock.AssertExpectations(t)
	})

	t.Run("fails with a repository error getting zones", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{}, nil)
		zonesRepositoryMock.On("GetAll", ctx).Return(nil, mockError)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.Error(t, err)
		require.Equal(t, mockError, err)
		require.Nil(t, res)

		pricesRepositoryMock.AssertExpectations(t)
		zonesRepositoryMock.AssertExpectations(t)
	})

	t.Run("fails with a repository error saving prices", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{testPrices}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(mockError)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.Error(t, err)
		require.Equal(t, mockError, err)
		require.Nil(t, res)

		pricesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
	})

	t.Run("fails with a provider error", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{testPrices}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return(nil, mockError)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		pricesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("provider does not return data", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{testPrices}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return([]domain.Prices{}, nil)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		pricesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("repository does not provide data", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{}, nil)
		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, todayTestDate).Return([]domain.Prices{testPricesFetch}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, ([]domain.Zone)(nil), tomorrowTestDate).Return([]domain.Prices{}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId}, res)

		pricesRepositoryMock.AssertExpectations(t)
		zonesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
	})

	t.Run("repository provides previous than today data and current hour <= 20", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 20, 0, 0, 0, time.Local)
		tomorrowDate := todayDate.AddDate(0, 0, 1)
		now = func() time.Time { return todayDate }

		yesterdayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.AddDate(0, 0, -1).Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{yesterdayPrices}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, todayDate).Return([]domain.Prices{testPricesFetch}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, ([]domain.Zone)(nil), tomorrowDate).Return([]domain.Prices{}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId}, res)

		pricesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
		zonesRepositoryMock.AssertNotCalled(t, "GetAll", ctx)
	})

	t.Run("repository provides today data and current hour <= 20", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 20, 0, 0, 0, time.Local)
		tomorrowDate := todayDate.AddDate(0, 0, 1)
		now = func() time.Time { return todayDate }

		todayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{todayPrices}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, ([]domain.Zone)(nil), todayDate).Return([]domain.Prices{}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, ([]domain.Zone)(nil), tomorrowDate).Return([]domain.Prices{}, nil)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		pricesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
		zonesRepositoryMock.AssertNotCalled(t, "GetAll", ctx)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("repository provides previous than today data and current hour > 20", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 21, 0, 0, 0, time.Local)
		tomorrowDate := todayDate.AddDate(0, 0, 1)
		now = func() time.Time { return todayDate }

		yesterdayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.AddDate(0, 0, -1).Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{yesterdayPrices}, nil)
		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, todayDate).Return([]domain.Prices{testPricesFetch}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, tomorrowDate).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId, testPricesFetchId}, res)

		pricesRepositoryMock.AssertExpectations(t)
		zonesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
	})

	t.Run("repository provides today data and current hour > 20", func(t *testing.T) {
		pricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 21, 0, 0, 0, time.Local)
		tomorrowDate := todayDate.AddDate(0, 0, 1)
		now = func() time.Time { return todayDate }

		todayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{todayPrices}, nil)
		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, ([]domain.Zone)(nil), todayDate).Return([]domain.Prices{}, nil)
		pricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, tomorrowDate).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(pricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId}, res)

		pricesRepositoryMock.AssertExpectations(t)
		zonesRepositoryMock.AssertExpectations(t)
		pricesProviderMock.AssertExpectations(t)
	})
}
