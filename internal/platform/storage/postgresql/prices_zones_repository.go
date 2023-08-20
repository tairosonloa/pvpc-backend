package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/charmbracelet/log"
	"github.com/huandu/go-sqlbuilder"

	pvpc "pvpc-backend/internal"
	"pvpc-backend/internal/errors"
)

// PricesZonesRepository is a PostgreSQL pvpc.PricesZonesRepository implementation.
type PricesZonesRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewPricesZonesRepository initializes a PostgreSQL-based implementation of pvpc.PricesZonesRepository.
func NewPricesZonesRepository(db *sql.DB, dbTimeout time.Duration) *PricesZonesRepository {
	return &PricesZonesRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *PricesZonesRepository) GetAll(ctx context.Context) ([]pvpc.PricesZone, error) {
	log.Debug("Getting all PricesZones from database")
	zoneStruct := sqlbuilder.NewStruct(new(zoneSchema))

	query, _ := zoneStruct.SelectFrom(zonesTableName).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctxTimeout, query)
	if err != nil {
		return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error querying PricesZone from database")
	}
	defer rows.Close()

	zones := make([]pvpc.PricesZone, 0, 5)
	for rows.Next() {
		var dbZone zoneSchema
		err := rows.Scan(zoneStruct.Addr(&dbZone)...)
		if err != nil {
			return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping PricesZone from database to schema")
		}

		zone, err := mapDbPricesZoneToDomain(dbZone)
		if err != nil {
			return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping PricesZone from schema to domain")
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

func (r *PricesZonesRepository) GetByID(ctx context.Context, id pvpc.PricesZoneID) (pvpc.PricesZone, error) {
	log.Debug("Getting PricesZone from database by ID", "id", id.String())
	zoneStruct := sqlbuilder.NewStruct(new(zoneSchema))

	qb := zoneStruct.SelectFrom(zonesTableName)
	query, args := qb.Where(qb.Equal("id", id.String())).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var dbZone zoneSchema
	err := row.Scan(zoneStruct.Addr(&dbZone)...)

	if err != nil {
		if err == sql.ErrNoRows {
			return pvpc.PricesZone{}, errors.NewDomainError(errors.PricesZoneNotFound, "PricesZone with ID %s not found", id.String())
		}
		return pvpc.PricesZone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping PricesZone from database to schema")
	}

	zone, err := mapDbPricesZoneToDomain(dbZone)
	if err != nil {
		return pvpc.PricesZone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping PricesZone from schema to domain")
	}

	return zone, nil
}
func (r *PricesZonesRepository) GetByExternalID(ctx context.Context, externalID string) (pvpc.PricesZone, error) {
	log.Debug("Getting PricesZone from database by externalID", "externalID", externalID)
	zoneStruct := sqlbuilder.NewStruct(new(zoneSchema))

	qb := zoneStruct.SelectFrom(zonesTableName)
	query, args := qb.Where(qb.Equal("external_id", externalID)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var dbZone zoneSchema
	err := row.Scan(zoneStruct.Addr(&dbZone)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return pvpc.PricesZone{}, errors.NewDomainError(errors.PricesZoneNotFound, "PricesZone with externalID %s not found", externalID)
		}
		return pvpc.PricesZone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping PricesZone from database to schema")
	}

	zone, err := mapDbPricesZoneToDomain(dbZone)
	if err != nil {
		return pvpc.PricesZone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping PricesZone from schema to domain")
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
