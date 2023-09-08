package postgresql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/huandu/go-sqlbuilder"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/domain/errors"
	"pvpc-backend/pkg/logger"
)

const (
	pricesTableName = "prices"
)

type pricesSchema struct {
	ID           string                 `db:"id"`
	Date         string                 `db:"date"`
	ZoneID       string                 `db:"zone_id"`
	HourlyPrices hourlyPriceSchemaSlice `db:"values"`
}

type hourlyPriceSchemaSlice []hourlyPriceSchema

type hourlyPriceSchema struct {
	Datetime string  `json:"datetime"`
	Price    float32 `json:"value"`
}

// Make the hourlyPriceSchemaSlice type implement the driver.Value interface.
// This method simply returns the JSON-encoded representation of the struct.
func (ps hourlyPriceSchemaSlice) Value() (driver.Value, error) {
	return json.Marshal(ps)
}

// Make the hourlyPriceSchemaSlice type implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (ps *hourlyPriceSchemaSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.NewDomainError(errors.PersistenceError, "sql.Scanner Scan() custom implementation: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ps)
}

// PricesRepository is a PostgreSQL domain.PricesRepository implementation.
type PricesRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewPricesRepository initializes a PostgreSQL-based implementation of domain.PricesRepository.
func NewPricesRepository(db *sql.DB, dbTimeout time.Duration) *PricesRepository {
	return &PricesRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

// Save implements the domain.PricesRepository interface.
func (r *PricesRepository) Save(ctx context.Context, prices []domain.Prices) error {
	logger.DebugContext(ctx, "Saving Prices into database")
	pricesSQL := sqlbuilder.NewStruct(new(pricesSchema))

	dbPrices := make([]interface{}, len(prices))

	for i, p := range prices {
		values := make([]hourlyPriceSchema, len(p.Values()))
		for j, v := range p.Values() {
			values[j] = hourlyPriceSchema{
				Datetime: v.Datetime().Format(time.RFC3339),
				Price:    v.Value(),
			}
		}

		dbPrices[i] = pricesSchema{
			ID:           p.ID().String(),
			Date:         p.Date().Format(time.RFC3339),
			ZoneID:       p.Zone().ID().String(),
			HourlyPrices: values,
		}
	}

	query, args := pricesSQL.InsertInto(pricesTableName, dbPrices...).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return errors.WrapIntoDomainError(err, errors.PersistenceError, "error trying to persist Prices into database")
	}

	return nil
}

// Query implements the domain.PricesRepository interface.
func (r *PricesRepository) Query(ctx context.Context, zoneID *domain.ZoneID, date *time.Time) ([]domain.Prices, error) {
	logger.DebugContext(ctx, "Getting all Zones from database")
	sb := sqlbuilder.NewSelectBuilder()
	pricesSQL := sqlbuilder.NewStruct(new(pricesSchema))

	query := pricesSQL.SelectFrom(pricesTableName)

	if date == nil {
		if zoneID == nil {
			query = pricesSQL.SelectFrom(pricesTableName).Distinct().SQL("ON zone_id").OrderBy("date").Desc()
		} else {
			query = query.Where(sb.Equal("zone_id", zoneID.String())).OrderBy("date").Desc().Limit(1)
		}
	} else {
		query = query.Where(sb.Equal("date", date.Format("2006-01-02")))
	}

	querySQL, args := query.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctxTimeout, querySQL, args...)
	if err != nil {
		return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error querying Prices from database")
	}
	defer rows.Close()

	prices := make([]domain.Prices, 0, 5)
	for rows.Next() {
		var dbPrices pricesSchema
		err := rows.Scan(pricesSQL.Addr(&dbPrices)...)
		if err != nil {
			return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Prices from database to schema")
		}

		domainPrices, err := mapPricesSchemaToDomain(dbPrices)
		if err != nil {
			return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Prices from schema to domain")
		}
		prices = append(prices, domainPrices)
	}

	return prices, nil
}

func mapPricesSchemaToDomain(priceSchema pricesSchema) (domain.Prices, error) {
	var hourlyPrices []domain.HourlyPriceDto

	for _, v := range priceSchema.HourlyPrices {
		hourlyPrice := domain.HourlyPriceDto{
			Datetime: v.Datetime,
			Value:    v.Price,
		}

		hourlyPrices = append(hourlyPrices, hourlyPrice)
	}

	return domain.NewPrices(domain.PricesDto{
		ID:     priceSchema.ID,
		Date:   priceSchema.Date,
		Zone:   domain.ZoneDto{},
		Values: hourlyPrices,
	})
}
