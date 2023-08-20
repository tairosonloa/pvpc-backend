package postgresql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/charmbracelet/log"
	"github.com/huandu/go-sqlbuilder"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/domain/errors"
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
	log.Debug("Saving Prices into database")
	pricesStruct := sqlbuilder.NewStruct(new(pricesSchema))

	dbPrices := make([]interface{}, len(prices))

	for i, p := range prices {
		values := make([]hourlyPriceSchema, len(p.Values()))
		for j, v := range p.Values() {
			values[j] = hourlyPriceSchema{
				Datetime: v.Datetime(),
				Price:    v.Value(),
			}
		}

		dbPrices[i] = pricesSchema{
			ID:           p.ID().String(),
			Date:         p.Date(),
			ZoneID:       p.Zone().ID().String(),
			HourlyPrices: values,
		}
	}

	query, args := pricesStruct.InsertInto(pricesTableName, dbPrices...).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return errors.WrapIntoDomainError(err, errors.PersistenceError, "error trying to persist Prices into database")
	}

	return nil
}
