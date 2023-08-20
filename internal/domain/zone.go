package domain

import (
	"context"
	"regexp"

	"pvpc-backend/internal/errors"
)

// ZoneDto is the DTO structure that represents a PVPC prices zone.
type ZoneDto struct {
	ID         string
	ExternalID string
	Name       string
}

// Zone represents a PVPC prices zone.
type Zone struct {
	id         ZoneID
	externalID string
	name       string
}

// ZoneID represents the PVPC prices zone unique identifier.
type ZoneID struct {
	value string
}

// NewZoneID instantiate the VO for ZoneID.
func NewZoneID(value string) (ZoneID, error) {
	// ID consist of 3 uppercase letters (A-Z) representing the zone.
	err := errors.NewDomainError(errors.InvalidZoneID, "invalid Zone ID: %s. It must be three capital letters", value)

	if len(value) != 3 {
		return ZoneID{}, err
	}

	if !regexp.MustCompile(`[A-Z]{3}`).MatchString(value) {
		return ZoneID{}, err
	}

	return ZoneID{
		value: value,
	}, nil
}

// String converts the ZoneID into string.
func (id ZoneID) String() string {
	return id.value
}

// ZonesRepository defines the expected behavior from a prices storage.
type ZonesRepository interface {
	// GetAll returns all the prices zones.
	GetAll(ctx context.Context) ([]Zone, error)
	// GetByID returns the prices zone with the given ID.
	GetByID(ctx context.Context, id ZoneID) (Zone, error)
	// GetByExternalID returns the prices zone with the given external ID.
	GetByExternalID(ctx context.Context, externalID string) (Zone, error)
}

// NewZone creates a new Zone struct.
func NewZone(zoneDto ZoneDto) (Zone, error) {
	idVO, err := NewZoneID(zoneDto.ID)
	if err != nil {
		return Zone{}, err
	}

	zone := Zone{
		id:         idVO,
		externalID: zoneDto.ExternalID,
		name:       zoneDto.Name,
	}

	return zone, nil
}

// ID returns the Zone unique identifier.
func (c Zone) ID() ZoneID {
	return c.id
}

// ExternalID returns the Zone external ID.
func (c Zone) ExternalID() string {
	return c.externalID
}

// Name returns the Zone Name.
func (c Zone) Name() string {
	return c.name
}
