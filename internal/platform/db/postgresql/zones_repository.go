package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pvpc "go-pvpc/internal"

	"github.com/huandu/go-sqlbuilder"
)

// PricesZoneRepository is a PostgreSQL pvpc.PricesZoneRepository implementation.
type PricesZoneRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewPricesZoneRepository initializes a PostgreSQL-based implementation of pvpc.PricesZoneRepository.
func NewPricesZoneRepository(db *sql.DB, dbTimeout time.Duration) *PricesZoneRepository {
	return &PricesZoneRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *PricesZoneRepository) GetAll(ctx context.Context) ([]pvpc.PricesZone, error) {
	zoneStruct := sqlbuilder.NewStruct(new(zoneSchema))

	query, _ := zoneStruct.SelectFrom(zonesTableName).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctxTimeout, query)
	if err != nil {
		return nil, fmt.Errorf("error querying prices zone from database: %v", err)
	}
	defer rows.Close()

	zones := make([]pvpc.PricesZone, 0, 5)
	for rows.Next() {
		var dbZone zoneSchema
		rows.Scan(zoneStruct.Addr(&dbZone)...)

		zone, err := mapDbPricesZoneToDomain(dbZone)
		if err != nil {
			return nil, fmt.Errorf("error mapping prices zone from database: %v", err)
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

func (r *PricesZoneRepository) GetByID(ctx context.Context, id pvpc.PricesZoneID) (pvpc.PricesZone, error) {
	zoneStruct := sqlbuilder.NewStruct(new(zoneSchema))

	qb := zoneStruct.SelectFrom(zonesTableName)
	query, args := qb.Where(qb.Equal("id", id.String())).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var dbZone zoneSchema
	row.Scan(zoneStruct.Addr(&dbZone)...)

	zone, err := mapDbPricesZoneToDomain(dbZone)
	if err != nil {
		return pvpc.PricesZone{}, fmt.Errorf("error mapping prices zone from database: %v", err)
	}

	return zone, nil
}
func (r *PricesZoneRepository) GetByExternalID(ctx context.Context, externalID string) (pvpc.PricesZone, error) {
	zoneStruct := sqlbuilder.NewStruct(new(zoneSchema))

	qb := zoneStruct.SelectFrom(zonesTableName)
	query, args := qb.Where(qb.Equal("external_id", externalID)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var dbZone zoneSchema
	row.Scan(zoneStruct.Addr(&dbZone)...)

	zone, err := mapDbPricesZoneToDomain(dbZone)
	if err != nil {
		return pvpc.PricesZone{}, fmt.Errorf("error mapping prices zone from database: %v", err)
	}

	return zone, nil
}

func mapDbPricesZoneToDomain(zoneSchema zoneSchema) (pvpc.PricesZone, error) {
	return pvpc.NewPricesZone(pvpc.PricesZoneDto{
		ID:         zoneSchema.ID,
		ExternalID: zoneSchema.ExternalID,
		Name:       zoneSchema.Name,
	})
}
