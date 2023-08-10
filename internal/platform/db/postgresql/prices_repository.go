package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"

	pvpc "go-pvpc/internal"
)

// PricesRepository is a PostgreSQL pvpc.PricesRepository implementation.
type PricesRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewPricesRepository initializes a PostgreSQL-based implementation of pvpc.PricesRepository.
func NewPricesRepository(db *sql.DB, dbTimeout time.Duration) *PricesRepository {
	return &PricesRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

// Save implements the pvpc.PricesRepository interface.
func (r *PricesRepository) Save(ctx context.Context, prices pvpc.Prices) error {
	pricesSQLStruct := sqlbuilder.NewStruct(new(sqlPrices))

	values := make([]sqlPrice, len(prices.Values()))
	for i, v := range prices.Values() {
		values[i] = sqlPrice{
			Datetime: v.Datetime(),
			Value:    v.Value(),
		}
	}

	query, args := pricesSQLStruct.InsertInto(sqlPricesTable, sqlPrices{
		ID:      prices.ID().String(),
		Date:    prices.Date(),
		GeoId:   prices.GeoId(),
		GeoName: prices.GeoName(),
		Values:  values,
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist prices on database: %v", err)
	}

	return nil
}
