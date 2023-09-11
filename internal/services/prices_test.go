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
	_, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: "2020-01-01T00:00:00Z", Values: []domain.HourlyPriceDto{testValuesDto}})
	require.NoError(t, err)
	testPricesFetchIdStr := "ZON-2023-01-01"
	testPricesFetchId, err := domain.NewPricesID(testPricesFetchIdStr)
	require.NoError(t, err)
	testPricesFetch, err := domain.NewPrices(domain.PricesDto{ID: testPricesFetchIdStr, Zone: testZoneDto, Date: "2023-01-01T00:00:00Z", Values: []domain.HourlyPriceDto{testValuesDto}})
	require.NoError(t, err)
	testZone, err := domain.NewZone(testZoneDto)
	require.NoError(t, err)
	loc, err := time.LoadLocation("Europe/Madrid")
	require.NoError(t, err)
	todayTestDate := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)
	tomorrowTestDate := todayTestDate.AddDate(0, 0, 1)
	now = func() time.Time { return todayTestDate }
	defer restoreNow(time.Now)

	t.Run("fails with a repository error getting zones", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		zonesRepositoryMock.On("GetAll", ctx).Return(nil, mockError)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.Error(t, err)
		require.Equal(t, mockError, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
	})

	t.Run("fails with a repository error getting prices", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return(nil, mockError)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.Error(t, err)
		require.Equal(t, mockError, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
	})

	t.Run("fails with a repository error saving prices", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(mockError)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.Error(t, err)
		require.Equal(t, mockError, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
	})

	t.Run("fails with a provider error", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		mockError := errors.NewDomainError(errors.PersistenceError, "mock-error")
		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return(nil, mockError)
		fallbackPricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return(nil, mockError)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
		fallbackPricesProviderMock.AssertExpectations(t)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("provider does not return data", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()
		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return([]domain.Prices{}, nil)
		fallbackPricesProviderMock.On("FetchPVPCPrices", ctx, mock.Anything, mock.Anything).Return([]domain.Prices{}, nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
		fallbackPricesProviderMock.AssertExpectations(t)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("repository does not provide data and current hour <= 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 20, 0, 0, 0, loc)
		now = func() time.Time { return todayDate }

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, todayTestDate).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId}, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
		mainPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", ctx, ([]domain.Zone)(nil), tomorrowTestDate)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repository does not provide data and current hour > 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 21, 0, 0, 0, loc)
		now = func() time.Time { return todayDate }

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, todayTestDate).Return([]domain.Prices{testPricesFetch}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, tomorrowTestDate).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId, testPricesFetchId}, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repository provides previous than today data and current hour <= 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 20, 0, 0, 0, loc)
		today := time.Date(todayDate.Year(), todayDate.Month(), todayDate.Day(), 0, 0, 0, 0, todayDate.Location())
		tomorrow := today.AddDate(0, 0, 1)
		now = func() time.Time { return todayDate }

		yesterdayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.AddDate(0, 0, -1).Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{yesterdayPrices}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, today).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId}, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
		mainPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", ctx, ([]domain.Zone)(nil), tomorrow)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repository provides previous than today data and current hour > 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 21, 0, 0, 0, loc)
		today := time.Date(todayDate.Year(), todayDate.Month(), todayDate.Day(), 0, 0, 0, 0, todayDate.Location())
		tomorrow := today.AddDate(0, 0, 1)
		now = func() time.Time { return todayDate }

		yesterdayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.AddDate(0, 0, -1).Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{yesterdayPrices}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, today).Return([]domain.Prices{testPricesFetch}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, tomorrow).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId, testPricesFetchId}, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repository provides today data and current hour <= 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 20, 0, 0, 0, loc)
		now = func() time.Time { return todayDate }

		todayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{todayPrices}, nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("repository provides today data and current hour > 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 21, 0, 0, 0, loc)
		today := time.Date(todayDate.Year(), todayDate.Month(), todayDate.Day(), 0, 0, 0, 0, todayDate.Location())
		tomorrow := today.AddDate(0, 0, 1)
		now = func() time.Time { return todayDate }

		todayPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{todayPrices}, nil)
		mainPricesProviderMock.On("FetchPVPCPrices", ctx, []domain.Zone{testZone}, tomorrow).Return([]domain.Prices{testPricesFetch}, nil)
		pricesRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Equal(t, []domain.PricesID{testPricesFetchId}, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertExpectations(t)
		mainPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", ctx, ([]domain.Zone)(nil), today)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repository provides tomorrow data and current hour <= 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 20, 0, 0, 0, loc)
		now = func() time.Time { return todayDate }

		tomorrowPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.AddDate(0, 0, 1).Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{tomorrowPrices}, nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("repository provides tomorrow data and current hour > 20", func(t *testing.T) {
		mainPricesProviderMock := new(mocks.PricesProvider)
		fallbackPricesProviderMock := new(mocks.PricesProvider)
		pricesRepositoryMock := new(mocks.PricesRepository)
		zonesRepositoryMock := new(mocks.ZonesRepository)
		ctx := context.Background()

		currentNow := now
		defer restoreNow(currentNow)
		todayDate := time.Date(2020, 1, 1, 21, 0, 0, 0, loc)
		now = func() time.Time { return todayDate }

		tomorrowPrices, err := domain.NewPrices(domain.PricesDto{ID: "ZON-2020-01-01", Zone: testZoneDto, Date: todayDate.AddDate(0, 0, 1).Format(time.RFC3339), Values: []domain.HourlyPriceDto{testValuesDto}})
		require.NoError(t, err)

		zonesRepositoryMock.On("GetAll", ctx).Return([]domain.Zone{testZone}, nil)
		pricesRepositoryMock.On("Query", ctx, (*domain.ZoneID)(nil), (*time.Time)(nil)).Return([]domain.Prices{tomorrowPrices}, nil)

		pricesService := NewPricesService(mainPricesProviderMock, fallbackPricesProviderMock, pricesRepositoryMock, zonesRepositoryMock)
		res, err := pricesService.FetchAndStorePricesFromREE(ctx)
		require.NoError(t, err)
		require.Nil(t, res)

		zonesRepositoryMock.AssertExpectations(t)
		pricesRepositoryMock.AssertExpectations(t)
		mainPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
		fallbackPricesProviderMock.AssertNotCalled(t, "FetchPVPCPrices", mock.Anything, mock.Anything, mock.Anything)
		pricesRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})
}
