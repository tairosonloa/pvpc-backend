package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/domain/errors"
	"pvpc-backend/pkg/logger"
)

const (
	zonesTableName = "zones"
)

type zoneSchema struct {
	ID         string `db:"id"`
	ExternalID string `db:"external_id"`
	Name       string `db:"name"`
}

// ZonesRepository is a PostgreSQL domain.ZonesRepository implementation.
type ZonesRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewZonesRepository initializes a PostgreSQL-based implementation of domain.ZonesRepository.
func NewZonesRepository(db *sql.DB, dbTimeout time.Duration) *ZonesRepository {
	return &ZonesRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *ZonesRepository) GetAll(ctx context.Context) ([]domain.Zone, error) {
	logger.DebugContext(ctx, "Getting all Zones from database")
	zoneSQL := sqlbuilder.NewStruct(new(zoneSchema))

	query, _ := sqlbuilder.WithFlavor(zoneSQL.SelectFrom(zonesTableName), sqlbuilder.PostgreSQL).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctxTimeout, query)
	if err != nil {
		return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error querying Zone from database")
	}
	defer rows.Close()

	zones := make([]domain.Zone, 0, 5)
	for rows.Next() {
		var dbZone zoneSchema
		err := rows.Scan(zoneSQL.Addr(&dbZone)...)
		if err != nil {
			return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Zone from database to schema")
		}

		zone, err := mapZoneSchemaToDomain(dbZone)
		if err != nil {
			return nil, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Zone from schema to domain")
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

func (r *ZonesRepository) GetByID(ctx context.Context, id domain.ZoneID) (domain.Zone, error) {
	logger.DebugContext(ctx, "Getting Zone from database by ID", "id", id.String())
	zoneSQL := sqlbuilder.NewStruct(new(zoneSchema))

	selectQB := zoneSQL.SelectFrom(zonesTableName)
	query, args := sqlbuilder.WithFlavor(selectQB.Where(selectQB.Equal("id", id.String())), sqlbuilder.PostgreSQL).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var dbZone zoneSchema
	err := row.Scan(zoneSQL.Addr(&dbZone)...)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Zone{}, errors.NewDomainError(errors.ZoneNotFound, "Zone with ID %s not found", id.String())
		}
		return domain.Zone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Zone from database to schema")
	}

	zone, err := mapZoneSchemaToDomain(dbZone)
	if err != nil {
		return domain.Zone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Zone from schema to domain")
	}

	return zone, nil
}
func (r *ZonesRepository) GetByExternalID(ctx context.Context, externalID string) (domain.Zone, error) {
	logger.DebugContext(ctx, "Getting Zone from database by externalID", "externalID", externalID)
	zoneSQL := sqlbuilder.NewStruct(new(zoneSchema))

	selectQB := zoneSQL.SelectFrom(zonesTableName)
	query, args := sqlbuilder.WithFlavor(selectQB.Where(selectQB.Equal("external_id", externalID)), sqlbuilder.PostgreSQL).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var dbZone zoneSchema
	err := row.Scan(zoneSQL.Addr(&dbZone)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Zone{}, errors.NewDomainError(errors.ZoneNotFound, "Zone with externalID %s not found", externalID)
		}
		return domain.Zone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Zone from database to schema")
	}

	zone, err := mapZoneSchemaToDomain(dbZone)
	if err != nil {
		return domain.Zone{}, errors.WrapIntoDomainError(err, errors.PersistenceError, "error mapping Zone from schema to domain")
	}

	return zone, nil
}

func mapZoneSchemaToDomain(zoneSchema zoneSchema) (domain.Zone, error) {
	return domain.NewZone(domain.ZoneDto{
		ID:         zoneSchema.ID,
		ExternalID: zoneSchema.ExternalID,
		Name:       zoneSchema.Name,
	})
}
