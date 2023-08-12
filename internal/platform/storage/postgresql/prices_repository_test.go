package postgresql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pvpc "go-pvpc/internal"
)

func Test_PricesRepository_Save(t *testing.T) {

	t.Run("when db returns error, repository returns error", func(t *testing.T) {
		id1, date1 := "ZON-2023-08-10", "2023-08-10"
		id2, date2 := "ZON-2023-08-10", "2023-08-10"
		zoneID, zoneExternalID, zoneName := "ZON", "123", "Test zone"
		datetime, value := "2023-08-10T00:00:00+02:00", float32(0.1234)

		prices1, err := pvpc.NewPrices(pvpc.PricesDto{
			ID:     id1,
			Date:   date1,
			Zone:   pvpc.PricesZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []pvpc.PriceDto{{Datetime: datetime, Value: float32(value)}, {Datetime: datetime, Value: float32(value)}},
		})
		require.NoError(t, err)

		prices2, err := pvpc.NewPrices(pvpc.PricesDto{
			ID:     id2,
			Date:   date2,
			Zone:   pvpc.PricesZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []pvpc.PriceDto{{Datetime: datetime, Value: value}, {Datetime: datetime, Value: value}},
		})
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		values := priceSchemaSlice{{Datetime: datetime, Price: value}, {Datetime: datetime, Price: value}}

		sqlMock.ExpectExec(
			"INSERT INTO prices (id, date, zone_id, values) VALUES (?, ?, ?, ?), (?, ?, ?, ?)").
			WithArgs(id1, date1, zoneID, values, id2, date2, zoneID, values).
			WillReturnError(errors.New("mock-error"))

		repo := NewPricesRepository(db, 1*time.Millisecond)

		err = repo.Save(context.Background(), []pvpc.Prices{prices1, prices2})

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("when everything goes OK, repository returns no error", func(t *testing.T) {
		id1, date1 := "ZON-2023-08-10", "2023-08-10"
		id2, date2 := "ZON-2023-08-10", "2023-08-10"
		zoneID, zoneExternalID, zoneName := "ZON", "123", "Test zone"
		datetime, value := "2023-08-10T00:00:00+02:00", float32(0.1234)

		prices1, err := pvpc.NewPrices(pvpc.PricesDto{
			ID:     id1,
			Date:   date1,
			Zone:   pvpc.PricesZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []pvpc.PriceDto{{Datetime: datetime, Value: float32(value)}, {Datetime: datetime, Value: float32(value)}},
		})
		require.NoError(t, err)

		prices2, err := pvpc.NewPrices(pvpc.PricesDto{
			ID:     id2,
			Date:   date2,
			Zone:   pvpc.PricesZoneDto{ID: zoneID, ExternalID: zoneExternalID, Name: zoneName},
			Values: []pvpc.PriceDto{{Datetime: datetime, Value: value}, {Datetime: datetime, Value: value}},
		})
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		values := priceSchemaSlice{{Datetime: datetime, Price: value}, {Datetime: datetime, Price: value}}

		sqlMock.ExpectExec(
			"INSERT INTO prices (id, date, zone_id, values) VALUES (?, ?, ?, ?), (?, ?, ?, ?)").
			WithArgs(id1, date1, zoneID, values, id2, date2, zoneID, values).
			WillReturnResult(sqlmock.NewResult(0, 2))

		repo := NewPricesRepository(db, 1*time.Millisecond)

		err = repo.Save(context.Background(), []pvpc.Prices{prices1, prices2})

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
	})

}
