package pvpc

import (
	"context"
	"regexp"

	"pvpc-backend/internal/errors"
)

// PricesZoneDto is the DTO structure that represents a PVPC prices zone.
type PricesZoneDto struct {
	ID         string
	ExternalID string
	Name       string
}

// PricesZone represents a PVPC prices zone.
type PricesZone struct {
	id         PricesZoneID
	externalID string
	name       string
}

// PricesZoneID represents the PVPC prices zone unique identifier.
type PricesZoneID struct {
	value string
}

// NewPricesZoneID instantiate the VO for PricesZoneID.
func NewPricesZoneID(value string) (PricesZoneID, error) {
	// ID consist of 3 uppercase letters (A-Z) representing the zone.
	err := errors.NewDomainError(errors.InvalidPricesZoneID, "invalid PricesZone ID: %s. It must be three capital letters", value)

	if len(value) != 3 {
		return PricesZoneID{}, err
	}

	if !regexp.MustCompile(`[A-Z]{3}`).MatchString(value) {
		return PricesZoneID{}, err
	}

	return PricesZoneID{
		value: value,
	}, nil
}

// String converts the PricesZoneID into string.
func (id PricesZoneID) String() string {
	return id.value
}

// PricesZonesRepository defines the expected behavior from a prices storage.
type PricesZonesRepository interface {
	// GetAll returns all the prices zones.
	GetAll(ctx context.Context) ([]PricesZone, error)
	// GetByID returns the prices zone with the given ID.
	GetByID(ctx context.Context, id PricesZoneID) (PricesZone, error)
	// GetByExternalID returns the prices zone with the given external ID.
	GetByExternalID(ctx context.Context, externalID string) (PricesZone, error)
}

// NewPricesZone creates a new PricesZone struct.
func NewPricesZone(pricesZoneDto PricesZoneDto) (PricesZone, error) {
	idVO, err := NewPricesZoneID(pricesZoneDto.ID)
	if err != nil {
		return PricesZone{}, err
	}

	pricesZone := PricesZone{
		id:         idVO,
		externalID: pricesZoneDto.ExternalID,
		name:       pricesZoneDto.Name,
	}

	return pricesZone, nil
}

// ID returns the PricesZone unique identifier.
func (c PricesZone) ID() PricesZoneID {
	return c.id
}

// ExternalID returns the PricesZone external ID.
func (c PricesZone) ExternalID() string {
	return c.externalID
}

// Name returns the PricesZone Name.
func (c PricesZone) Name() string {
	return c.name
}
