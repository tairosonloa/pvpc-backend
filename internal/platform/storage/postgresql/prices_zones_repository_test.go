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
	dErrors "go-pvpc/internal/errors"
)

func Test_ZonesRepository_GetAll(t *testing.T) {

	t.Run("when db returns error, repository returns error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones").
			WillReturnError(errors.New("mock-error"))

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		_, err = repo.GetAll(context.Background())

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("when there are zones in the database, returns a slice of pvpc.PricesZones", func(t *testing.T) {
		id1, externalID1, name1 := "ZON", "123", "Test zone 1"
		id2, externalID2, name2 := "ABC", "456", "Test zone 2"

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "external_id", "name"}).
			AddRow(id1, externalID1, name1).
			AddRow(id2, externalID2, name2)

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones").
			WillReturnRows(rows)

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		result, err := repo.GetAll(context.Background())
		require.NoError(t, err)

		expected1, err := pvpc.NewPricesZone(pvpc.PricesZoneDto{ID: id1, ExternalID: externalID1, Name: name1})
		require.NoError(t, err)

		expected2, err := pvpc.NewPricesZone(pvpc.PricesZoneDto{ID: id2, ExternalID: externalID2, Name: name2})
		require.NoError(t, err)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expected1, result[0])
		assert.Equal(t, expected2, result[1])
	})

	t.Run("when there are NOT zones in the database, returns an empty slice of pvpc.PricesZones", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "external_id", "name"})

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones").
			WillReturnRows(rows)

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		result, err := repo.GetAll(context.Background())
		require.NoError(t, err)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})
}

func Test_ZonesRepository_GetByID(t *testing.T) {

	t.Run("when db returns error, repository returns error", func(t *testing.T) {
		zoneIDString := "ZON"
		zoneID, err := pvpc.NewPricesZoneID(zoneIDString)
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones WHERE id = ?").
			WithArgs(zoneIDString).
			WillReturnError(errors.New("mock-error"))

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		_, err = repo.GetByID(context.Background(), zoneID)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("when db returns no error and zone is found, repository returns a pvpc.PricesZone", func(t *testing.T) {
		id, externalID, name := "ZON", "123", "Test zone"
		zoneID, err := pvpc.NewPricesZoneID(id)
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "external_id", "name"}).
			AddRow(id, externalID, name)

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones WHERE id = ?").
			WithArgs(id).
			WillReturnRows(rows)

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		result, err := repo.GetByID(context.Background(), zoneID)
		require.NoError(t, err)

		expected, err := pvpc.NewPricesZone(pvpc.PricesZoneDto{ID: id, ExternalID: externalID, Name: name})
		require.NoError(t, err)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, expected, result)

	})

	t.Run("when db returns no error but zone is NOT found, repository returns an empty pvpc.PricesZone and error", func(t *testing.T) {
		id := "ZON"
		zoneID, err := pvpc.NewPricesZoneID(id)
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "external_id", "name"})

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones WHERE id = ?").
			WithArgs(id).
			WillReturnRows(rows)

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		result, err := repo.GetByID(context.Background(), zoneID)

		assert.Error(t, err)
		assert.Equal(t, dErrors.PricesZoneNotFound, dErrors.Code(err))
		assert.Equal(t, pvpc.PricesZone{}, result)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func Test_ZonesRepository_GetByExternalID(t *testing.T) {

	t.Run("when db returns error, repository returns error", func(t *testing.T) {
		zoneExternalID := "123"

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones WHERE external_id = ?").
			WithArgs(zoneExternalID).
			WillReturnError(errors.New("mock-error"))

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		_, err = repo.GetByExternalID(context.Background(), zoneExternalID)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("when db returns no error, repository returns a pvpc.PricesZone", func(t *testing.T) {
		id, externalID, name := "ZON", "123", "Test zone"

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "external_id", "name"}).
			AddRow(id, externalID, name)

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones WHERE external_id = ?").
			WithArgs(externalID).
			WillReturnRows(rows)

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		result, err := repo.GetByExternalID(context.Background(), externalID)
		require.NoError(t, err)

		expected, err := pvpc.NewPricesZone(pvpc.PricesZoneDto{ID: id, ExternalID: externalID, Name: name})
		require.NoError(t, err)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("when db returns no error but zone is NOT found, repository returns an empty pvpc.PricesZone and error", func(t *testing.T) {
		externalID := "123"
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "external_id", "name"})

		sqlMock.ExpectQuery("SELECT zones.id, zones.external_id, zones.name FROM zones WHERE external_id = ?").
			WithArgs(externalID).
			WillReturnRows(rows)

		repo := NewPricesZonesRepository(db, 1*time.Millisecond)

		result, err := repo.GetByExternalID(context.Background(), externalID)

		assert.Error(t, err)
		assert.Equal(t, dErrors.PricesZoneNotFound, dErrors.Code(err))
		assert.Equal(t, pvpc.PricesZone{}, result)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}
