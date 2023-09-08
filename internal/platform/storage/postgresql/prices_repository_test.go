package postgresql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"pvpc-backend/internal/domain"
)

func Test_PricesRepository_Save(t *testing.T) {

	t.Run("when db returns error, repository returns error", func(t *testing.T) {
		id1, date1 := "ZON-2023-08-10", "2023-08-10T00:00:00+02:00"
		id2, date2 := "ZON-2023-08-10", "2023-08-10T00:00:00+02:00"
		zoneID, zoneExternalID, zoneName := "ZON", "123", "Test zone"
		datetime, value := "2023-08-10T00:00:00+02:00", float32(0.1234)

		prices1, err := domain.NewPrices(domain.PricesDto{
			ID:     id1,
			Date:   date1,
			Zone:   domain.ZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []domain.HourlyPriceDto{{Datetime: datetime, Value: float32(value)}, {Datetime: datetime, Value: float32(value)}},
		})
		require.NoError(t, err)

		prices2, err := domain.NewPrices(domain.PricesDto{
			ID:     id2,
			Date:   date2,
			Zone:   domain.ZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []domain.HourlyPriceDto{{Datetime: datetime, Value: value}, {Datetime: datetime, Value: value}},
		})
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		values := hourlyPriceSchemaSlice{{Datetime: datetime, Price: value}, {Datetime: datetime, Price: value}}

		sqlMock.ExpectExec(
			"INSERT INTO prices (id, date, zone_id, values) VALUES (?, ?, ?, ?), (?, ?, ?, ?)").
			WithArgs(id1, date1, zoneID, values, id2, date2, zoneID, values).
			WillReturnError(errors.New("mock-error"))

		repo := NewPricesRepository(db, 1*time.Millisecond)

		err = repo.Save(context.Background(), []domain.Prices{prices1, prices2})

		require.NoError(t, sqlMock.ExpectationsWereMet())
		require.Error(t, err)
	})

	t.Run("when everything goes OK, repository returns no error", func(t *testing.T) {
		id1, date1 := "ZON-2023-08-10", "2023-08-10T00:00:00+02:00"
		id2, date2 := "ZON-2023-08-10", "2023-08-10T00:00:00+02:00"
		zoneID, zoneExternalID, zoneName := "ZON", "123", "Test zone"
		datetime, value := "2023-08-10T00:00:00+02:00", float32(0.1234)

		prices1, err := domain.NewPrices(domain.PricesDto{
			ID:     id1,
			Date:   date1,
			Zone:   domain.ZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []domain.HourlyPriceDto{{Datetime: datetime, Value: float32(value)}, {Datetime: datetime, Value: float32(value)}},
		})
		require.NoError(t, err)

		prices2, err := domain.NewPrices(domain.PricesDto{
			ID:     id2,
			Date:   date2,
			Zone:   domain.ZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []domain.HourlyPriceDto{{Datetime: datetime, Value: value}, {Datetime: datetime, Value: value}},
		})
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		values := hourlyPriceSchemaSlice{{Datetime: datetime, Price: value}, {Datetime: datetime, Price: value}}

		sqlMock.ExpectExec(
			"INSERT INTO prices (id, date, zone_id, values) VALUES (?, ?, ?, ?), (?, ?, ?, ?)").
			WithArgs(id1, date1, zoneID, values, id2, date2, zoneID, values).
			WillReturnResult(sqlmock.NewResult(0, 2))

		repo := NewPricesRepository(db, 1*time.Millisecond)

		err = repo.Save(context.Background(), []domain.Prices{prices1, prices2})

		require.NoError(t, sqlMock.ExpectationsWereMet())
		require.NoError(t, err)
	})

}

func Test_PricesRepository_Query(t *testing.T) {

	t.Run("when db returns error, repository returns error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		sqlMock.ExpectQuery(
			"SELECT DISTINCT ON (zone_id) id, date, zone_id, values FROM prices ORDER BY date DESC").
			WillReturnError(errors.New("mock-error"))

		repo := NewPricesRepository(db, 1*time.Millisecond)

		_, err = repo.Query(context.Background(), nil, nil)

		require.NoError(t, sqlMock.ExpectationsWereMet())
		require.Error(t, err)
	})

	t.Run("queries by zone ID", func(t *testing.T) {
		date := "2023-08-10T00:00:00+02:00"

		id, err := domain.NewPricesID("ZON-2023-08-10")
		require.NoError(t, err)

		zoneID, err := domain.NewZoneID("ZON")
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "date", "zone_id", "values"}).
			AddRow(id.String(), date, zoneID.String(), hourlyPriceSchemaSlice{{Datetime: date, Price: float32(0.1234)}})

		sqlMock.ExpectQuery(
			"SELECT prices.id, prices.date, prices.zone_id, prices.values FROM prices WHERE zone_id = ZON ORDER BY date DESC LIMIT 1").
			WillReturnRows(rows)

		repo := NewPricesRepository(db, 1*time.Millisecond)

		result, err := repo.Query(context.Background(), &zoneID, nil)
		require.NoError(t, err)

		prices, err := domain.NewPrices(domain.PricesDto{ID: id.String(), Date: date, Zone: domain.ZoneDto{ID: zoneID.String()}, Values: []domain.HourlyPriceDto{{Datetime: date, Value: float32(0.1234)}}})
		require.NoError(t, err)

		require.Equal(t, []domain.Prices{prices}, result)

		require.NoError(t, sqlMock.ExpectationsWereMet())
		require.NoError(t, err)
	})

}
