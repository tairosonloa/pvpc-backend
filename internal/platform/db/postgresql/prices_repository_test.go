package postgresql

import (
	"context"
	"errors"
	pvpc "go-pvpc/internal"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CourseRepository_Save_RepositoryError(t *testing.T) {
	id1, date1 := "ZON-2023-08-10", "2023-08-10"
	id2, date2 := "ZON-2023-08-10", "2023-08-10"
	zoneId, zoneExternalId, zoneName := "ZON", "123", "Test zone"
	datetime, value := "2023-08-10T00:00:00+02:00", float32(0.1234)

	prices1, err := pvpc.NewPrices(pvpc.PricesDto{
		ID:     id1,
		Date:   date1,
		Zone:   pvpc.PricesZoneDto{ID: zoneId, ExternalId: zoneExternalId, Name: zoneName},
		Values: []pvpc.PriceDto{{Datetime: datetime, Value: float32(value)}, {Datetime: datetime, Value: float32(value)}},
	})
	require.NoError(t, err)

	prices2, err := pvpc.NewPrices(pvpc.PricesDto{
		ID:     id2,
		Date:   date2,
		Zone:   pvpc.PricesZoneDto{ID: zoneId, ExternalId: zoneExternalId, Name: zoneName},
		Values: []pvpc.PriceDto{{Datetime: datetime, Value: value}, {Datetime: datetime, Value: value}},
	})
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	values := priceSchemaSlice{{Datetime: datetime, Price: value}, {Datetime: datetime, Price: value}}

	sqlMock.ExpectExec(
		"INSERT INTO prices (id, date, zone_id, values) VALUES (?, ?, ?, ?), (?, ?, ?, ?)").
		WithArgs(id1, date1, zoneId, values, id2, date2, zoneId, values).
		WillReturnError(errors.New("test-error"))

	repo := NewPricesRepository(db, 1*time.Millisecond)

	err = repo.Save(context.Background(), []pvpc.Prices{prices1, prices2})

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

// func Test_CourseRepository_Save_Succeed(t *testing.T) {
// 	courseID, courseName, courseDuration := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"

// 	course, err := pvpc.NewCourse(courseID, courseName, courseDuration)
// 	require.NoError(t, err)

// 	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	require.NoError(t, err)

// 	sqlMock.ExpectExec(
// 		"INSERT INTO courses (id, name, duration) VALUES (?, ?, ?)").
// 		WithArgs(courseID, courseName, courseDuration).
// 		WillReturnResult(sqlmock.NewResult(0, 1))

// 	repo := NewCourseRepository(db, 1*time.Millisecond)

// 	err = repo.Save(context.Background(), course)

// 	assert.NoError(t, sqlMock.ExpectationsWereMet())
// 	assert.NoError(t, err)
// }
